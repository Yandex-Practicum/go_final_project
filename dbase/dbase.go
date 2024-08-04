package dbase

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func getDBFilePath() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(filepath.Dir(appPath), "../scheduler.db")
	}
	return dbFile
}

func InitializeDB() (*sql.DB, error) {
	dbFilePath := getDBFilePath()
	_, err := os.Stat(dbFilePath)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	if install {
		createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			date    CHAR(8) NOT NULL DEFAULT "",
			title   VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat  VARCHAR(128) NOT NULL DEFAULT ""
		);`

		_, err = db.Exec(createTableSQL)
		if err != nil {
			log.Println("Error creating table", err)
			return nil, err
		}

		createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
		_, err = db.Exec(createIndexSQL)
		if err != nil {
			log.Println("Error creating index", err)
			return nil, err
		}
	}
	return db, nil
}
