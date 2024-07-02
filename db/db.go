package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() {
	dbName := "./db/scheduler.db"
	dbDir := "/Users/dkotov/Desktop/Practicum/go_final_project/"
	dbtodo := os.Getenv("TODO_DBFILE")

	dbFile := filepath.Join(filepath.Dir(dbDir), dbName)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	} else {
		fmt.Println("Database already exists at:", dbFile)
	}
	if dbtodo != "" {
		fmt.Println("Database created successfully at TODO_DBFILE:", dbtodo)
	} else {
		if install {
			db, err := sql.Open("sqlite3", dbFile)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()
			createTabSQL := `CREATE TABLE IF NOT EXISTS scheduler (
				"id"	INTEGER,
				"date"	TEXT NOT NULL,
				"title"	TEXT NOT NULL,
				"comment"	TEXT,
				"repeat"	TEXT NOT NULL DEFAULT "",
				CHECK(length("repeat") <= 128),
				CHECK(length("title") > 0),
				PRIMARY KEY("id" AUTOINCREMENT)
			);
			CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`
			_, err = db.Exec(createTabSQL)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Database created successfully at:", dbFile)

		}
	}
}
