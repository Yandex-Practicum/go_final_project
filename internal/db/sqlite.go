package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName = "scheduler.db"
)

func New() *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal("init db", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("ping db", err)
	}

	return db

}
