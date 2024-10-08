package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func setupDB() (*sql.DB, error) {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("ошибка директории: %v", err)
		}
		dbPath = filepath.Join(cwd, "scheduler.db")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к БД: %v", err)
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		createTableQuery := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128)
		);
		CREATE INDEX idx_date ON scheduler(date);
		`
		if _, err := db.Exec(createTableQuery); err != nil {
			db.Close()
			return nil, fmt.Errorf("ошибка при создании таблицы: %v", err)
		}
	}

	return db, nil
}
