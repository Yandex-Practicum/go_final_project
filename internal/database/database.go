package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitializeDatabase() *sql.DB {
	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "todo-server", "scheduler.db")
	environment := os.Getenv("TODO_DBFILE")
	if len(environment) > 0 {
		dbFile = environment
	}
	_, err = os.Stat(filepath.Join(dbFile))

	if os.IsNotExist(err) {
		db, err := sql.Open("sqlite", "scheduler.db")
		if err != nil {
			log.Fatal(err)
		}

		query := `
		CREATE TABLE IF NOT EXISTS "scheduler" (
			"id" INTEGER NOT NULL UNIQUE,
			"date" CHAR(8) NOT NULL DEFAULT "",
			"title" VARCHAR(128) NOT NULL DEFAULT "",
			"comment" TEXT NOT NULL DEFAULT "",
			"repeat" VARCHAR(128) NOT NULL DEFAULT "",
			PRIMARY KEY("id")
		);`

		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}

		query = `CREATE INDEX IF NOT EXISTS "scheduler_index_date" ON "scheduler" ("date");`
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}

		return db
	} else {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
