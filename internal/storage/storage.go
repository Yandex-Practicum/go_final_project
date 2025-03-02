package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func OpenStorage(storagePath string) (*Storage, error) {

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("database open error: %w\n", err)
	}

	if pingErr := db.Ping(); pingErr != nil {
		return nil, fmt.Errorf("database connection error: %w\n", pingErr)
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
		return nil, fmt.Errorf("database create error: %w\n", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date on scheduler(date);`)
	if err != nil {
		return nil, fmt.Errorf("index create error: %w\n", err)
	}
	return &Storage{db: db}, nil
}
