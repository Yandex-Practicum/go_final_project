package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go_final_project-main/internal/model"
	"net/http"
	"time"
)

func GetTasksHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		var tasks []model.Task
		var err error

		// Начинаем формировать SQL-запрос
		query := `SELECT * FROM scheduler`
		var args []interface{}

		// Проверка наличия параметра поиска
		if search != "" {
			if _, err := time.Parse("02.01.2006", search); err == nil {
				// Форматируем дату для SQL
				parsedDate, _ := time.Parse("02.01.2006", search)
				date := parsedDate.Format("20060102")

				// Запрос для выборки по дате
				query += ` WHERE date = ?`
				args = append(args, date)
			} else {
				// Запрос для выборки по заголовку или комментарию
				query += ` WHERE title LIKE ? OR comment LIKE ?`
				args = append(args, "%"+search+"%", "%"+search+"%")
			}
		}

		query += ` ORDER BY date LIMIT 50`

		// Выполняем запрос
		err = db.Select(&tasks, query, args...)
		if err != nil {
			http.Error(w, `{"error": "Ошибка при выборке задач"}`, http.StatusInternalServerError)
			return
		}

		// Если задач нет, возвращаем пустой слайс
		if tasks == nil {
			tasks = []model.Task{}
		}

		// Формируем ответ
		response := model.TasksResponse{
			Tasks: make([]model.TaskResponse, len(tasks)),
		}

		for i, task := range tasks {
			response.Tasks[i] = model.TaskResponse{
				ID:      fmt.Sprint(task.ID), // Преобразуем ID в строку
				Date:    task.Date,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error": "Ошибка при формировании ответа"}`, http.StatusInternalServerError)
		}
	}
}
