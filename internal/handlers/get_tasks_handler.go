package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project-main/internal/model"
	"go_final_project-main/internal/nextdate"

	"github.com/jmoiron/sqlx"
)

const orderByLimit = `50`

func GetTasksHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		var tasks []model.Task
		var err error

		// Начинаем формировать SQL-запрос
		query := `SELECT id, date, title, comment, repeat FROM scheduler`
		var args []interface{}

		const SearchDateFormat = "02.01.2006"
		// Проверка наличия параметра поиска
		if search != "" {
			if _, err := time.Parse(SearchDateFormat, search); err == nil {
				// Форматируем дату для SQL
				parsedDate, _ := time.Parse(SearchDateFormat, search)
				date := parsedDate.Format(nextdate.DateFormat)

				// Запрос для выборки по дате
				query += ` WHERE date = ?`
				args = append(args, date)
			} else {
				// Запрос для выборки по заголовку или комментарию
				query += ` WHERE title LIKE ? OR comment LIKE ?`
				args = append(args, "%"+search+"%", "%"+search+"%")
			}
		}

		query += ` ORDER BY date LIMIT ` + orderByLimit

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
			Tasks: make([]model.Task, len(tasks)),
		}

		for i, task := range tasks {
			response.Tasks[i] = model.Task{
				ID:      task.ID,
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
