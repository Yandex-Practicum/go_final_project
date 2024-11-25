package repository

import (
	"database/sql"

	"go_final_project/internal/task"
)

func (rep *Repository) GetTaskByID(id string) (task.Task, error) {
	var t task.Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id= ?`
	row := rep.db.QueryRow(query, id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil || err == sql.ErrNoRows {
		return t, err
	}

	return t, nil
}
