package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"todo_restapi/internal/dto"
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

func (s *Storage) GetAllTasks() ([]models.Task, error) {
	
}
