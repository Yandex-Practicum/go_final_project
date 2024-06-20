package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func dbExists() bool {
	dbFile := os.Getenv("TODO_DBFILE")
	_, err := os.Stat(dbFile)
	var exists bool
	if err == nil {
		exists = true
	}
	return exists

}

func installDB() {
	dbFile := os.Getenv("TODO_DBFILE")
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	installQuery, err := os.ReadFile("backend/install.sql")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec(string(installQuery))
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()
}
