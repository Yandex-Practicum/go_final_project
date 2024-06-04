package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

var response struct {
	Error string `json:"error,omitempty"`
}

func CreateDB(DBFile string) {
	db, err := sql.Open("sqlite3", DBFile)

	if err != nil {
		return
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			"id" INTEGER PRIMARY KEY,
			"date" CHAR(8) NOT NULL,
			"title" VARCHAR(255) NOT NULL,
			"comment" TEXT,
			"repeat" VARCHAR(128) NOT NULL
			)
		`)
	if err != nil {
		return
	}
	_, err = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_date ON sсheduler (date)
		`)
	if err != nil {
		return
	}
}

func main() {
	if debugNextDate {
		line := []string{"20240409", "m 31", "20240531"}
		now, err := time.Parse(dataFormat, "20240126")
		if err != nil {
			return
		}
		date := line[0]
		repeat := line[1]
		repeatDate, err := NextDate(now, date, repeat)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(repeatDate)
		return
	}
	ServerStart()
}
