package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		_, err = db.Exec(`CREATE TABLE scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT
        )`)
		if err != nil {
			return nil, err
		}
		_, err = db.Exec("CREATE INDEX idx_date ON scheduler (date)")
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
