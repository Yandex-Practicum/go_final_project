package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func GetDBFilePath() (string, error) {
	dbFilePath := os.Getenv("TODO_DBFILE")
	if dbFilePath == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		dbFilePath = filepath.Join(currentDir, "scheduler.db")
	}
	return dbFilePath, nil
}

func InitDB(dbFilePath string) {
	var err error
	DB, err = sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date TEXT NOT NULL,
  title TEXT NOT NULL,
  comment TEXT,
  repeat TEXT CHECK(length(repeat) <= 128)
 );`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`)
	if err != nil {
		log.Fatal(err)
	}
}
