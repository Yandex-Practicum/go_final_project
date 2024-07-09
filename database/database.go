package database

import (
	"database/sql"
	"os"
	"path/filepath"
)

func InitDb() (*sql.DB, error) {
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

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return nil, err
	}

	if !install {
		queryCreate, err := db.Query(`
			CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title TEXT NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX date_index ON scheduler (date);
		`)
		if err != nil {
			return nil, err
		}
		defer queryCreate.Close()
	}

	return db, nil
}
