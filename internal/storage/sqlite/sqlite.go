package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"todo-list/internal/storage"

	_ "modernc.org/sqlite"
)

const dbFilePath = "internal/storage/sqlite/scheduler.db"

type Storage struct {
	db *sql.DB
}

func NewStorage(log *slog.Logger) (*Storage, error) {

	dbPath, err := storage.DBFilePath(dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get database file path:%w", err)
	}

	var install bool
	_, err = os.Stat(dbPath)
	if err != nil {
		install = true
		log.Debug("Getting ready for creating database")
	} else {
		log.Debug("Database file is found")
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with SQLite database: %w", err)
	}

	if install {
		stmt, err := db.Prepare(`CREATE TABLE scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
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

		log.Debug("Database file is created")
	}

	return &Storage{db: db}, nil
}
