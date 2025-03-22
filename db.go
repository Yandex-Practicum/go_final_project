package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const defaultDBFile = "scheduler.db"

func getDBPath() string {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = defaultDBFile
	}
	return dbPath
}

func connectDB() (*sql.DB, error) {
	dbPath := getDBPath()

	_, err := os.Stat(dbPath)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД %s: %w", dbPath, err)
	}

	if install {
		if err := createTable(db); err != nil {
			db.Close()
			return nil, fmt.Errorf("ошибка создания таблицы в БД %s: %w", dbPath, err)
		}
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	createTableSql := `
	CREATE TABLE IF NOT EXISTS scheduler (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    date TEXT NOT NULL,
	    title TEXT NOT NULL,
	    comment TEXT,
	    repeat TEXT CHECK(length(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`
	_, err := db.Exec(createTableSql)
	return err
}
