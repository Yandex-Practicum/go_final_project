package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DBFile = "./scheduler.db"

func InitDatabase() (*sql.DB, error) {
	dbfile := DBFile
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbfile = envFile
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Printf("База данных и таблица успешно созданы.")
	return db, nil
}
