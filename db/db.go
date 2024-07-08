package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() (*sqlx.DB, error) {
	dbFileName := "scheduler.db"
	dbFile := filepath.Join(".", dbFileName)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	} else {
		fmt.Println("Database already exists at:", dbFile)
	}

	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath != "" {
		fmt.Println("Database created successfully at TODO_DBFILE:", dbPath)
		dbFile = dbPath
	} else if install {
		db, err := sqlx.Connect("sqlite3", dbFile)
		if err != nil {
			return nil, err
		}

		createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			return nil, err
		}

		createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`
		_, err = db.Exec(createIndexSQL)
		if err != nil {
			return nil, err
		}

		fmt.Println("Database created successfully at:", dbFile)
	}

	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}
