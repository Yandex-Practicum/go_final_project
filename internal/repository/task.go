package repository

import (
	"go_final_project/config"
	"go_final_project/internal/models"
)

func AddTaskToDB(task *models.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := config.GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
