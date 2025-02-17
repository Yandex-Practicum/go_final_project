package repository

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"go_final_project/internal/model"
)

type TaskRepository interface {
	AddTask(task model.Task) (int64, error)
	GetTasksByDate(date string, limit int) ([]model.Task, error)
	GetTasksBySearch(search string, limit int) ([]model.Task, error)
	GetAllTasks(limit int) ([]model.Task, error)
	GetTaskByID(id string) (model.Task, error)
	UpdateTask(task model.Task) error
	DeleteTask(id int) error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) AddTask(task model.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *taskRepository) GetTasksByDate(date string, limit int) ([]model.Task, error) {
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE date = ? 
		ORDER BY date 
		LIMIT ?
	`
	rows, err := r.db.Query(query, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func (r *taskRepository) GetTasksBySearch(search string, limit int) ([]model.Task, error) {
	searchPattern := "%" + search + "%"
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE (title LIKE ? OR comment LIKE ?) 
		ORDER BY date 
		LIMIT ?
	`
	rows, err := r.db.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func (r *taskRepository) GetAllTasks(limit int) ([]model.Task, error) {
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		ORDER BY date 
		LIMIT ?
	`
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

func (r *taskRepository) GetTaskByID(id string) (model.Task, error) {
	var task model.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, errors.New("задача не найдена")
		}
		return task, err
	}
	return task, nil
}

func (r *taskRepository) UpdateTask(task model.Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	_, err := r.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	return err
}

func (r *taskRepository) DeleteTask(id int) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func scanTasks(rows *sql.Rows) ([]model.Task, error) {
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
