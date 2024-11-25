package repository

import (
	"database/sql"

	"go_final_project/internal/task"
)

func (r *Repository) GetTasks() ([]task.Task, error) {
	var tasks []task.Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler  ORDER BY date ASC  LIMIT 20`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil || err == sql.ErrNoRows {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
