package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Читаем параметр search из строки запроса
		search := r.URL.Query().Get("search")
		var query string
		var args []interface{}
		limit := 50 // Ограничение на количество записей

		if search == "" {
			// Запрос для получения всех задач без фильтрации
			query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
			args = append(args, limit)
		} else {
			// Проверка: соответствует ли параметр формату даты
			if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
				// Дату 02.01.2006 преобразуем в 20060102
				formattedDate := parsedDate.Format("20060102")
				query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
				args = append(args, formattedDate, limit)
			} else {
				// Ищем по заголовку и комментарию с использованием LIKE
				likePattern := "%" + search + "%"
				query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
				args = append(args, likePattern, likePattern, limit)
			}
		}

		// Выполняем запрос
		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при выполнении запроса: %v"}`, err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Считываем результаты запроса
		tasks := []Task{}
		for rows.Next() {
			var task Task
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка при чтении данных: %v"}`, err), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		// Проверяем ошибки после цикла чтения
		if err := rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при чтении строк: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Если задач нет, возвращаем пустой список
		if len(tasks) == 0 {
			tasks = []Task{}
		}

		// Отправляем JSON-ответ
		response := map[string][]Task{"tasks": tasks}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при кодировании ответа"}`, http.StatusInternalServerError)
			return
		}

		w.Write(jsonResponse)
	}
}
