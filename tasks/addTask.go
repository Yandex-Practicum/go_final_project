package tasks

import "database/sql"

func AddTask(task Task) (int64, error) {
	var id int64
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, _ = res.LastInsertId()
	}
	return id, err
}
