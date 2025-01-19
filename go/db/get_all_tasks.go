package db

import (
	"database/sql"
	"final-project-bronner/go/models"
	"fmt"
)

func (db *DB) GetAllTasks() ([]models.Task, error) {
	limit := 20
	tasks := make([]models.Task, 0, 10)
	var res *sql.Rows
	var err error

	res, err = db.db.Query(
		"SELECT * FROM scheduler ORDER BY date LIMIT :limit",
		sql.Named("limit", limit),
	)

	if err != nil {
		return tasks, fmt.Errorf("ошибка при получении задач из бд: %w", err)
	}
	for res.Next() {
		var t models.Task
		err := res.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return tasks, fmt.Errorf("ошибка при мапинге задач из бд: %w", err)
		}
		if t.Title != "" {
			tasks = append(tasks, t)
		}
	}
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при мапинге задач из бд: %w", err)
	}
	return tasks, nil
}
