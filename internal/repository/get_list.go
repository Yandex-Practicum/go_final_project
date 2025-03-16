package repository

import (
	"fmt"

	"github.com/askoren1/go_final_project/internal/models"
)

// Функция для получения списка задач из базы данных
func (r *Repository) GetList() ([]models.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 30;` // формируем SQL-запрос
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task //Инициализация списка задач

	for rows.Next() { //Итерация по результатам запроса
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return tasks, err //Возврат результата
}
