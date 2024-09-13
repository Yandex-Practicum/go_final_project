package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(log *slog.Logger) (*Storage, error) {

	appPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to perform os.Executable(): %w", err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
		log.Info("Getting ready for creating database")
	} else {
		log.Info("Database file is found")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with SQLite database: %w", err)
	}

	if install {
		stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY AOUTOENCREMENT,)
			date CHAR(8),
			title VARCHAR(256),
			comment VARCHAR(512),
			repeat VARCHAR(128))
		`)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare query for creating table: %w", err)
		}

		_, err = stmt.Exec()
		if err != nil {
			return nil, fmt.Errorf("failed to create table in database: %w", err)
		}
	}

	return &Storage{db: db}, nil
}
