package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
	"todo-list/internal/storage"
	"todo-list/internal/tasks"

	_ "github.com/mattn/go-sqlite3"
)

const dbFilePath = "database/scheduler.db"

type Storage struct {
	db *sql.DB
}

func NewStorage(log *slog.Logger) (*Storage, error) {

	dbPath, err := storage.DBFilePath(dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get database file path:%w", err)
	}

	var install bool
	_, err = os.Stat(dbPath)
	if err != nil {
		install = true
		log.Debug("Getting ready for creating database")
	} else {
		log.Debug("Database file is found")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with SQLite database: %w", err)
	}

	if install {
		stmt, err := db.Prepare(`CREATE TABLE scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8),
			title VARCHAR(256),
			comment VARCHAR(512),
			repeat VARCHAR(128))
		`)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare query for creating table: %w", err)
		}

		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("failed to create table in database: %w", err)
		}

		log.Debug("Database file is created")
	}

	return &Storage{db: db}, nil
}

func (storage Storage) AddTask(task *tasks.Task) (int, error) {

	query := `INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)`

	result, err := storage.db.Exec(query, sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, fmt.Errorf("failed to insert into scheduler: %w", err)
	}

	ind, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last id inserted into scheduler: %w", err)
	}

	return int(ind), nil
}

func (storage Storage) GetTasks() ([]tasks.Task, error) {

	result := make([]tasks.Task, 0)
	_ = result

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY scheduler.date LIMIT 30"
	rows, err := storage.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks list from scheduler: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		task := tasks.Task{}

		err = rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row result from scheduler to task.Task: %w", err)
		}

		result = append(result, task)
	}

	return result, nil
}

func (storage Storage) GetTask(taskId string) (*tasks.Task, error) {

	id, err := strconv.Atoi(taskId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert taskId to int: %w", err)
	}

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id"
	row := storage.db.QueryRow(query, sql.Named("id", id))

	result := tasks.Task{}
	err = row.Scan(&result.Id, &result.Date, &result.Title, &result.Comment, &result.Repeat)
	if err != nil {
		return nil, fmt.Errorf("failed to read task from query result: %w", err)
	}

	return &result, nil
}

func (storage Storage) UpdateTask(task *tasks.Task) error {

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"

	taskId, err := strconv.Atoi(task.Id)
	if err != nil {
		return fmt.Errorf("failed convert task.id to int: %w", err)
	}

	res, err := storage.db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", taskId))
	if err != nil {
		return fmt.Errorf("failed to update task in database: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read database response: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id \"%d\" is not found", taskId)
	}

	return nil
}

func (storage Storage) MarkAsDone(taskId string) error {

	task, err := storage.GetTask(taskId)
	if err != nil {
		return fmt.Errorf("failed to get task by id: %w", err)
	}

	if task.Repeat == "" {
		err = storage.DeleteTask(taskId)
		if err != nil {
			return fmt.Errorf("failed to delete task by id: %w", err)
		}
	} else {
		now := time.Now()
		nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("failed to get the next task date: %w", err)
		}
		task.Date = nextDate
		err = storage.UpdateTask(task)
		if err != nil {
			return fmt.Errorf("failed to update task with new date: %w", err)
		}
	}

	return nil
}

func (storage Storage) DeleteTask(taskId string) error {

	id, err := strconv.Atoi(taskId)
	if err != nil {
		return fmt.Errorf("failed to convert task id to int: %w", err)
	}

	query := "DELETE FROM scheduler WHERE id = :id"
	res, err := storage.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("failed to delete task by id from database: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read query result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task by id \"%d\" is not found", id)
	}

	return nil
}
