package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const (
	createTable = `CREATE TABLE scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT "",
title VARCHAR(32) NOT NULL DEFAULT "",
comment TEXT NOT NULL DEFAULT "",
repeat VARCHAR(128) NOT NULL DEFAULT ""
);`
	createIndex = `CREATE INDEX idx_date ON scheduler (date);`
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error load env: %w", err)
	}

	dbFile := "scheduler.db"
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}

	_, err := os.Stat(dbFile)
	needInit := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("Could not open DB: %w", err)
	}

	if needInit {
		if _, err := db.Exec(createTable); err != nil {
			return nil, fmt.Errorf("Could not create table: %w", err)
		}
		if _, err := db.Exec(createIndex); err != nil {
			return nil, fmt.Errorf("Could not create index: %w", err)
		}
	}

	return db, nil
}
