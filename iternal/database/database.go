package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() (*sql.DB, error) {
	//appPath, err := os.Executable()
	//if err != nil {
	//log.Fatal(err)
	//}
	dbFile := "project.db"
	log.Println(dbFile)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	if install {
		query := ` 
		CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT CHECK(LENGTH(repeat) <= 128)
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`
		_, err = db.Exec(query)
		if err != nil {
			return nil, err
		}
		log.Println("База данных создана")
	}
	return db, nil
}
