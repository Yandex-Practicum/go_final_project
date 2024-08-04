package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbx *sqlx.DB

func initDb() {
	dbFile := getDbFilePath()
	_, err := os.Stat(dbFile)
	needCreate := err != nil
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbx, err = sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	if needCreate {
		installDb(db)
		log.Println("database created in: " + dbFile)
	} else {
		log.Println("existing database found in: " + dbFile)
	}

}

func getDbFilePath() string {
	return "scheduler.db"
}

func installDb(db *sql.DB) {
	sqlStmt := `CREATE TABLE scheduler (id INTEGER NOT NULL PRIMARY KEY, date VARCHAR(10), title NVARCHAR(255), comment NVARCHAR(10000), repeat VARCHAR(10));`
	_, err := dbx.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}

const HUNDRED = 100

func getAllTasks() ([]Task, error) {

	var tasks []Task

	err := dbx.Select(&tasks, "SELECT * FROM scheduler ORDER BY date ASC LIMIT HUNDRED")
	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

func loadTaskById(id int64) (*Task, error) {

	var task Task

	err := dbx.Get(&task, "SELECT * FROM scheduler where id = ?", id)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func insertTask(date string, title, comment, repeat string) (int, error) {

	sqlStmt := `INSERT INTO scheduler (date, title , comment , repeat) VALUES (?, ?, ?, ?)`
	r, err := dbx.Exec(sqlStmt, date, title, comment, repeat)
	if err != nil {
		return -1, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil

}

func deleteTaskById(id int) error {

	sqlStmt := `DELETE FROM scheduler WHERE id = ?`
	_, err := dbx.Exec(sqlStmt, id)
	if err != nil {
		return err
	}

	return nil
}

func updateTask(id int, date string, title, comment, repeat string) error {

	sqlStmt := `
	UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ? 
	WHERE id = ?`
	_, err := dbx.Exec(sqlStmt, date, title, comment, repeat, id)
	if err != nil {
		return err
	}

	return nil
}
