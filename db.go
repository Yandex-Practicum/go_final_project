package main

import (
	"database/sql"
	"log"
	"os"

	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	_, err = os.Stat(dbFile)
	install := os.IsNotExist(err)
	log.Println("Using database file:", dbFile)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		err = createTables(db)
		if err != nil {
			return nil, err
		}
		log.Println("Database created")
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT CHECK(length(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
	`
	_, err := db.Exec(query)
	return err
}
