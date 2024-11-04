package api

import (
	"database/sql"
	"encoding/json"
	"go_final_project/internal/config"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

// Response представляет ответ от сервера
type Response struct {
	Tasks []config.Task `json:"tasks"`
	Error string        `json:"error,omitempty"`
}

// GetUpcomingTasks возвращает список ближайших задач.
func GetUpcomingTasks(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Получаем параметр поиска из URL
		search := r.URL.Query().Get("search")

		// SQL запрос для получения задач
		var rows *sql.Rows
		var err error

		// Проверка на наличие даты в формате 02.01.2006
		if search != "" {
			if _, err := time.Parse("02.01.2006", search); err == nil {
				// Преобразуем строку в time.Time
				parsedDate, parseErr := time.Parse("02.01.2006", search) // Преобразуем строку в time.Time
				if parseErr != nil {
					respondWithError(w, http.StatusBadRequest, "неправильный формат даты")
					return
				}
				date := parsedDate.Format("20060102") // Получаем дату в нужном формате

				// Выполняем SQL-запрос по дате
				rows, err = db.Query(`
                    SELECT id, date, title, comment, repeat
                    FROM scheduler 
                    WHERE date = ?
                    ORDER BY date LIMIT 50`, date)
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
			} else {
				// Форматируем строку поиска для LIKE
				searchPattern := "%" + search + "%"
				rows, err = db.Query(`
                    SELECT id, date, title, comment, repeat
                    FROM scheduler 
                    WHERE title LIKE ? OR comment LIKE ?
                    ORDER BY date LIMIT 50`, searchPattern, searchPattern)
				if err != nil {
					respondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		} else {
			// Если параметр поиска пустой, возвращаем все задачи
			rows, err = db.Query(`
                SELECT id, date, title, comment, repeat
                FROM scheduler 
                ORDER BY date LIMIT 50`)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		defer rows.Close()

		var tasks []config.Task

		for rows.Next() {
			var task config.Task
			if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			tasks = append(tasks, task)
		}

		if tasks == nil {
			tasks = []config.Task{}
		}

		response := Response{Tasks: tasks}
		json.NewEncoder(w).Encode(response)
	}
}

// respondWithError отправляет JSON-ответ с ошибкой.
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{Error: message})
}
