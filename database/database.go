package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const defDBFile = "./scheduler.db"

func InitDatabase() (*sql.DB, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defDBFile
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %v", err)
	}

	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='scheduler'").Scan(&tableExists)
	if err != nil {
		return nil, fmt.Errorf("can't check if table exists: %v", err)
	}

	if tableExists == 0 {
		todolist := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128)
		);
		CREATE INDEX idx_date ON scheduler (date);
		`

		_, err = db.Exec(todolist)
		if err != nil {
			return nil, fmt.Errorf("can't create table or index: %v", err)
		}
		fmt.Printf("Database file created: %s\n", dbFile)
		fmt.Println("Table created")
	} else {
		fmt.Printf("Using existing database: %s\n", dbFile)
	}

	return db, nil
}
