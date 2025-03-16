package database

import (
	"fmt"

	"go_final_project/tasks"
)

func (d *Database) AddTask(t tasks.Task) (int, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat)
			 VALUES (?, ?, ?, ?)`

	// Выполнение запроса на добавление записи в таблицу
	res, err := d.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to add task: %w", err)
	}

	// Получение идентификатора добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve task ID: %w", err)
	}

	return int(id), nil
}
