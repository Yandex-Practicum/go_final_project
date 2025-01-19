package db

import (
	"database/sql"
	"final-project-bronner/go/models"
	"fmt"
)

func (db *DB) GetTask(id int) (models.Task, error) {
	var t models.Task
	res := db.db.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	)
	err := res.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return t, fmt.Errorf("ошибка при получении задачи: %w", err)
	}
	return t, nil
}
