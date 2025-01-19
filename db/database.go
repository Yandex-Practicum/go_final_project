package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func GetDatabasePath() string {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
		dbPath = filepath.Join(workingDir, "scheduler.db")
	}
	return dbPath
}

func SetupDatabase(dbFile string) error {
	_, err := os.Stat(dbFile)
	var dbCreated bool
	if err != nil && os.IsNotExist(err) {
		dbCreated = true
		log.Println("Database file not found, creating an empty database file.")
		file, err := os.Create(dbFile)
		if err != nil {
			return fmt.Errorf("failed to create database file: %v", err)
		}
		file.Close()
	}

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if dbCreated {
		err = createTable(db)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenDB() (*sql.DB, error) {
	dbPath := GetDatabasePath()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	log.Println("Creating table 'scheduler'...")
	query := `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT CHECK(length(repeat) <= 128),
        completed INTEGER DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	log.Println("Table 'scheduler' created successfully.")
	return nil
}
