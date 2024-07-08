package db

import (
	"github.com/jmoiron/sqlx"
)

func GetTaskByID(s *sqlx.DB, id string) (Task, error) {
	var task Task
	err := s.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat,
	)
	return task, err
}
