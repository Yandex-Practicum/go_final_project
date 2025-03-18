package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"todo_restapi/internal/dto"
	"todo_restapi/internal/myfunctions"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func OpenStorage(storagePath string) (*Storage, error) {

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("database open error: %w", err)
	}

	if pingErr := db.Ping(); pingErr != nil {
		return nil, fmt.Errorf("database connection error: %w", pingErr)
	} else {
		fmt.Println("Connected to database!")
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date CHAR(8) NOT NULL DEFAULT '',
    	title TEXT NOT NULL DEFAULT '',
    	comment TEXT NOT NULL DEFAULT '',
    	repeat VARCHAR(128) NOT NULL DEFAULT '');
	`)
	if err != nil {
		return nil, fmt.Errorf("database create error: %w", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date on scheduler(date);`)
	if err != nil {
		return nil, fmt.Errorf("index create error: %w", err)
	}
	return NewStorage(db), nil
}
func (s *Storage) AddTask(task models.Task) (int64, error) {

	statement, err := s.db.Prepare("INSERT INTO scheduler(date, title, comment, repeat) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("statement prepration error: %w", err)
	}

	defer statement.Close()

	result, err := statement.Exec(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("statement execution error: %w", err)
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting ID error: %w", err)
	}
	return taskID, nil
}

func (s *Storage) GetTasks() ([]models.Task, error) {

	output := make([]models.Task, 0, 10)

	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10")
	if err != nil {
		return output, fmt.Errorf("row query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {

		var getTasks models.Task

		err := rows.Scan(&getTasks.ID, &getTasks.Date, &getTasks.Title, &getTasks.Comment, &getTasks.Repeat)
		if err != nil {
			return output, fmt.Errorf("row scan error: %w\n", err)
		}

		getTasks.ID = fmt.Sprint(getTasks.ID)

		output = append(output, getTasks)
	}

	if err := rows.Err(); err != nil {
		return output, fmt.Errorf("row iteration error: %w", err)
	}
	return output, nil
}

func (s *Storage) GetOneTask(id string) (models.Task, error) {

	var getTask models.Task

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return getTask, fmt.Errorf("parse ID error: %w", err)
	}

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?", parsedID)

	err = row.Scan(&getTask.ID, &getTask.Date, &getTask.Title, &getTask.Comment, &getTask.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return getTask, fmt.Errorf("task with id %v not found", id)
	} else if err != nil {
		return getTask, fmt.Errorf("scan error: %w", err)
	}

	return getTask, nil
}

func (s *Storage) EditTask(task models.Task) error {

	result, err := s.db.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %v not found", task.ID)
	}

	return nil
}

func (s *Storage) DeleteTask(id string) error {

	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("parse ID error: %w", err)
	}

	result, err := s.db.Exec("DELETE FROM scheduler WHERE id=?", parsedID)
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %v not found", parsedID)
	}

	return nil
}

func (s *Storage) SearchTasks(searchQuery string) ([]models.Task, error) {

	var query string
	var arguments []interface{}
	output := make([]models.Task, 0, 10)

	date, err := myfunctions.IsDate(searchQuery)
	if err == nil {
		query = "SELECT * FROM scheduler WHERE date=? LIMIT 10"
		arguments = append(arguments, date)
	} else {
		query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 10"
		searchPattern := "%" + searchQuery + "%"
		arguments = append(arguments, searchPattern, searchPattern)
	}

	rows, err := s.db.Query(query, arguments...)
	if err != nil {
		return output, fmt.Errorf("row query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {

		var getTasks models.Task

		err := rows.Scan(&getTasks.ID, &getTasks.Date, &getTasks.Title, &getTasks.Comment, &getTasks.Repeat)
		if err != nil {
			return output, fmt.Errorf("row scan error: %w\n", err)
		}

		getTasks.ID = fmt.Sprint(getTasks.ID)

		output = append(output, getTasks)
	}

	if err := rows.Err(); err != nil {
		return output, fmt.Errorf("row iteration error: %w", err)
	}
	return output, nil

}
