package db

import (
	"database/sql"
	"log"

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
/*func AddTask(task handler.Task) (int64, error) {
	// добавляем задачу в базу данных
	var t Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	res, err := db.Exec("INSERT INTO tasks (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}*/
