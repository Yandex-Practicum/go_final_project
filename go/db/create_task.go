package db

import (
	"database/sql"
	"final-project-bronner/go/models"
	"fmt"
)

func (db *DB) AddTask(t models.Task) (int, error) {
	res, err := db.db.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)

	if err != nil {
		return 0, fmt.Errorf("ошибка при вставке задачи: %w", err)
	}
	idLast, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при вставке задачи: %w", err)
	}
	return int(idLast), nil
}
