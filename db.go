package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func getDBFilePath() string {
	dbFilePath := os.Getenv("TODO_DBFILE")
	return dbFilePath
}

func connectDB() (*sql.DB, error) {
	dbFilePath := getDBFilePath()
	if dbFilePath == "" {
		dbFilePath = "scheduler.db"
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	createStmt := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			title TEXT,
			comment TEXT,
			repeat TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	_, err := db.Exec(createStmt)
	return err
}

func maina() {
	dbFilePath := os.Getenv("TODO_DBFILE")
	if dbFilePath != "" {
		fmt.Printf("Using custom database file path: %s\n", dbFilePath)
	}

	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database and table created successfully.")
}
