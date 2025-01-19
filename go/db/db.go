package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func getDbFilePath() (string, bool) {
	appPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	filename := os.Getenv("TODO_DBFILE")

	if filename == "" {
		log.Fatal("Необходимо указать путь к файлу БД в TODO_DBFILE")
	}

	dbFile := filepath.Join(filepath.Dir(appPath), filename)

	_, err = os.Stat(dbFile)

	return dbFile, err == nil
}

func CreateConnectionDB() *sql.DB {
	dbFile, dbFileAlreadyExists := getDbFilePath()

	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatal(err)
	}

	if !dbFileAlreadyExists {
		db.Exec(`CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT
		)`)
	}

	return db
}

func GetDB(connect *sql.DB) *DB {
	return &DB{
		db: connect,
	}
}
