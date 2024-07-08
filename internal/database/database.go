package database

import (
	"database/sql"
	"log"
)

func CreateTable(db *sql.DB) {
	createSchedulerTableSQL := `CREATE TABLE scheduler (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,	
		"date" TEXT, 	
		"title" TEXT,
		"comment" TEXT,
		"repeat" VARCHAR(128) NULL		
	  );`

	statement, err := db.Prepare(createSchedulerTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}
