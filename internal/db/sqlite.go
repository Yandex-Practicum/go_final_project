package db

import (
	"database/sql"
	"fmt"
	"log"

	"pwd/internal/handler"

	_ "modernc.org/sqlite"
)

func New() *sql.DB {

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal("init db", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("ping bd", err)
	}
	return db
}

// добавляем задачу в базу данных
func AddTask(db *sql.DB, task *handler.Task) (int, error) {
	query := "INSERT INTO tasks (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	return int(id), nil
}
