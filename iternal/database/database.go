package database

import (
	"Go/iternal/services"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "project.db"
)

func CreateDB() (*sql.DB, error) {
	//appPath, err := os.Executable()
	//if err != nil {
	//log.Fatal(err)
	//}
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		return nil, err
	}
	defer db.Close()

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

func PutTaskInDB(task services.Task) (int64, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetCountOfTasks() (int, error) {
	var count int64
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT count(*) FROM scheduler")
	_ = row.Scan(&count)

	return int(count), nil
}

func GetAllTasks() (*sql.Rows, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM 'scheduler' ORDER BY date")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
