package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() *sql.DB {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func addTask(task Task) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func getTasks() ([]Task, error) {
	rows, err := db.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func doneTask(id int) error {
	var date, repeat string
	err := db.QueryRow(`
		SELECT date, repeat FROM scheduler WHERE id = ?
	`, id).Scan(&date, &repeat)
	if err != nil {
		return err
	}

	if repeat == "" {

		_, err = db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	} else {

		newDate, err := NextDate(time.Now(), date, repeat)
		if err != nil {
			return err
		}
		_, err = db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, newDate, id)
	}
	return err
}
