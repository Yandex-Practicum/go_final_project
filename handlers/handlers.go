package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	models "go_final_project/models"
	tasks "go_final_project/tasks"
	"net/http"
	"time"
)

// HandleNextDate обработчик для API запроса /api/nextdate
func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid now format", http.StatusBadRequest)
		return
	}

	nextDate, err := tasks.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

// обработчик для API запроса /api/task
func HandleTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		switch r.Method {
		case http.MethodPost:
			var task models.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				fmt.Println("Ошибка десериализации JSON: ", err)
				http.Error(w, "Ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			if task.Title == "" {
				fmt.Println("Не указан заголовок задачи")
				http.Error(w, "Не указан заголовок задачи", http.StatusBadRequest)
				return
			}

			now := time.Now()
			if task.Date == "" {
				task.Date = now.Format("20060102")
			} else {
				date, err := time.Parse("20060102", task.Date)
				if err != nil {
					fmt.Println("Дата представлена в неправильном формате")
					http.Error(w, "Дата представлена в неправильном формате", http.StatusBadRequest)
					return
				}

				if date.Before(now) {
					if task.Repeat == "" {
						task.Date = now.Format("20060102")
					} else {
						nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
						if err != nil {
							fmt.Println("Ошибка вычисления следующей даты: ", err)
							http.Error(w, "Ошибка вычисления следующей даты: "+err.Error(), http.StatusBadRequest)
							return
						}
						task.Date = nextDate
					}
				}
			}

			if task.Repeat != "" {
				if err := tasks.ValidateRepeatRule(task.Repeat); err != nil {
					fmt.Println("Правило повторения указано в неправильном формате")
					http.Error(w, "Правило повторения указано в неправильном формате", http.StatusBadRequest)
					return
				}

				nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					fmt.Println("Ошибка вычисления следующей даты: ", err)
					http.Error(w, "Ошибка вычисления следующей даты: "+err.Error(), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
				fmt.Println("Cледующая даты: ", nextDate)
			}

			id, err := tasks.AddTask(db, task)
			if err != nil {
				fmt.Println("Ошибка добавления задачи: ", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			response := map[string]interface{}{
				"id": id,
			}
			json.NewEncoder(w).Encode(response)
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
				return
			}

			var task models.Task
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
			err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
				} else {
					fmt.Println("Ошибка выполнения запроса: ", err)
					http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
				}
				return
			}
			json.NewEncoder(w).Encode(task)
		case http.MethodPut:
			var task models.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				fmt.Println("Ошибка десериализации JSON: ", err)
				http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
				return
			}

			if task.ID == "" {
				fmt.Println("Не указан идентификатор задачи")
				http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
				return
			}

			if task.Title == "" {
				fmt.Println("Не указан заголовок задачи")
				http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
				return
			}

			now := time.Now()
			if task.Date == "" {
				task.Date = now.Format("20060102")
			} else {
				date, err := time.Parse("20060102", task.Date)
				if err != nil {
					fmt.Println("Дата представлена в неправильном формате")
					http.Error(w, `{"error": "Дата представлена в неправильном формате"}`, http.StatusBadRequest)
					return
				}

				if date.Before(now) {
					if task.Repeat == "" {
						task.Date = now.Format("20060102")
					} else {
						nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
						if err != nil {
							fmt.Println("Ошибка вычисления следующей даты: ", err)
							http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusBadRequest)
							return
						}
						task.Date = nextDate
					}
				}
			}

			if task.Repeat != "" {
				if err := tasks.ValidateRepeatRule(task.Repeat); err != nil {
					fmt.Println("Правило повторения указано в неправильном формате")
					http.Error(w, `{"error": "Правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
					return
				}
			}

			query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
			res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
			if err != nil {
				fmt.Println("Ошибка выполнения запроса: ", err)
				http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
				return
			}

			rowsAffected, err := res.RowsAffected()
			if err != nil {
				fmt.Println("Ошибка получения результата запроса: ", err)
				http.Error(w, `{"error": "Ошибка получения результата запроса"}`, http.StatusInternalServerError)
				return
			}

			if rowsAffected == 0 {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{})
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
				return
			}

			deleteQuery := `DELETE FROM scheduler WHERE id = ?`
			res, err := db.Exec(deleteQuery, id)
			if err != nil {
				fmt.Println("Ошибка выполнения запроса: ", err)
				http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
				return
			}

			rowsAffected, err := res.RowsAffected()
			if err != nil {
				fmt.Println("Ошибка получения результата запроса: ", err)
				http.Error(w, `{"error": "Ошибка получения результата запроса"}`, http.StatusInternalServerError)
				return
			}

			if rowsAffected == 0 {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{})
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func HandleGetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50`)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса: ", err)
			http.Error(w, "Ошибка выполнения запроса: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		tasks := []models.Task{} // Инициализируем пустой слайс

		for rows.Next() {
			var task models.Task
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				fmt.Println("Ошибка чтения строки: ", err)
				http.Error(w, "Ошибка чтения строки: "+err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			fmt.Println("Ошибка обработки результата: ", err)
			http.Error(w, "Ошибка обработки результата: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"tasks": tasks,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func HandleMarkTaskDone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		var task models.Task
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				fmt.Println("Ошибка выполнения запроса: ", err)
				http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
			}
			return
		}

		if task.Repeat == "" {
			// Удаляем одноразовую задачу
			deleteQuery := `DELETE FROM scheduler WHERE id = ?`
			_, err := db.Exec(deleteQuery, id)
			if err != nil {
				fmt.Println("Ошибка удаления задачи: ", err)
				http.Error(w, `{"error": "Ошибка удаления задачи"}`, http.StatusInternalServerError)
				return
			}
		} else {
			// Рассчитываем следующую дату для периодической задачи
			now := time.Now()
			nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				fmt.Println("Ошибка вычисления следующей даты: ", err)
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}

			// Обновляем дату задачи
			updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
			_, err = db.Exec(updateQuery, nextDate, id)
			if err != nil {
				fmt.Println("Ошибка обновления задачи: ", err)
				http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
				return
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
