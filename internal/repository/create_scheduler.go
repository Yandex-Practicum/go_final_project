package repository

import (
	"log"
)

func (rep *Repository) CreateScheduler() {
	query := `CREATE TABLE IF NOT EXISTS scheduler(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT(128))`
	_, err := rep.db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
