package db

import (
	"database/sql"
	"fmt"
	"os"

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
	dbFile := "scheduler.db"
	if envFile := os.Getenv("DBFILE"); envFile != "" {
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
