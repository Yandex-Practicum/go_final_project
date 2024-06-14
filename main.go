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

func (t *Task) MakeValid() (bool, string, *Task) {
	newTask := t

	todayDateStr := time.Now().Format(dataFormat)
	if len(t.Date) == 0 || t.Date == "" {
		t.Date = todayDateStr
		return true, "", newTask
	}

	date, err := time.Parse(dataFormat, t.Date)
	if err != nil {
		return false, err.Error(), newTask
	}

	if date.Before(time.Now()) {
		ruleIsSet := !(len(t.Repeat) == 0 || t.Repeat == "")
		if !ruleIsSet {
			t.Date = todayDateStr
			return true, todayDateStr, newTask
		}
		if ruleIsSet {
			newDateStr, err := NextDate(time.Now(), t.Date, t.Repeat)

			if err != nil {
				return false, err.Error(), newTask
			}
			t.Date = newDateStr
			return true, "", newTask
		}
	}
	return true, "", newTask
}

func (t *Task) IsValid() (bool, string) {
	if len(t.Title) == 0 {
		return false, "no title"
	}
	return true, ""
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
