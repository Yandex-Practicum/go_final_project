package repository

import (
	"database/sql"
	"strconv"

	"github.com/askoren1/go_final_project/internal/models"
)

// Функция для получения задачи по ее идентификатору (ID) из базы данных
func (r *Repository) GetTaskByID(id string) (*models.Task, error) {

	var task models.Task
	var taskID int64

	row := r.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id) //формируем SQL-запрос
	err := row.Scan(&taskID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	task.ID = strconv.FormatInt(taskID, 10) //Преобразование идентификатора в строковое представление

	return &task, nil //Возврат результата
}
