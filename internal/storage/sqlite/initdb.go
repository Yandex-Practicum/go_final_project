package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	log.Printf("using database file: %s", dbFile)

	db, err := sql.Open("sqlite3", storagePath+"?_fk=1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date VARCHAR(128) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128)
		)
	`); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date)
	`); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
