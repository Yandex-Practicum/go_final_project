package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func initDb() {
	dbFile := getDbFilePath()
	_, err := os.Stat(dbFile)
	if err != nil {
		installDb(dbFile)
	} else {
		fmt.Println("existing database found in: " + dbFile)
	}

}

func getDbFilePath() string {
	appPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	return dbFile
}

func installDb(dbFile string) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table scheduler (id integer not null primary key, date text, title text, comment text, repeat text);
	
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	fmt.Println("database created in: " + dbFile)
}

func getAllTasks() ([]Task, error) {

	db, err := sqlx.Connect("sqlite3", getDbFilePath())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var tasks []Task

	err = db.Select(&tasks, "SELECT * FROM scheduler ORDER BY date ASC")

	if err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil

}
func loadTaskById(id int64) (*Task, error) {

	db, err := sqlx.Connect("sqlite3", getDbFilePath())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var task Task

	err = db.Get(&task, "SELECT * FROM scheduler where id = ?", id)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func insertTask(date string, title, comment, repeat string) (int, error) {

	db, err := sql.Open("sqlite3", getDbFilePath())
	if err != nil {
		return -1, err
	}
	defer db.Close()

	sqlStmt := `insert into scheduler (date, title , comment , repeat) values (?, ?, ?, ?)`
	r, err := db.Exec(sqlStmt, date, title, comment, repeat)
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
	db, err := sql.Open("sqlite3", getDbFilePath())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `delete from scheduler where id = ?`
	_, err = db.Exec(sqlStmt, id)
	if err != nil {
		return err
	}
	return nil

}

func updateTask(id int, date string, title, comment, repeat string) error {

	db, err := sql.Open("sqlite3", getDbFilePath())
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `update scheduler
	set date = ?, title = ?, comment = ?, repeat = ? 
	where id = ?`
	_, err = db.Exec(sqlStmt, date, title, comment, repeat, id)
	if err != nil {
		return err
	}

	return nil

}
