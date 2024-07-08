package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func AddTask(s *sqlx.DB, task Task) (int64, error) {
	var id int64

	res, err := s.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, _ = res.LastInsertId()
	}
	return id, err
}
