package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"todo-list/internal/storage"
	"todo-list/internal/tasks"

	_ "github.com/mattn/go-sqlite3"
)

const dbFilePath = "internal/storage/sqlite/scheduler.db"

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

func (storage Storage) AddTask(task tasks.Task) (int, error) {

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
