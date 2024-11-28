package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"my_education/go/go_final_project/internal/logic"
	"net/http"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что запрос выполнен методом POST
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"только метод POST поддерживается"}`, http.StatusMethodNotAllowed)
			return
		}

		var task Task

		// Десериализуем JSON-запрос
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		// Проверяем обязательное поле title
		if task.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		// Проверяем и парсим поле date
		var taskDate time.Time
		if task.Date == "" {
			taskDate = time.Now().Truncate(24 * time.Hour)
		} else {
			taskDate, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
				return
			}
		}

		// Проверяем, если дата меньше сегодняшнего дня
		now := time.Now().Truncate(24 * time.Hour)
		if taskDate.Before(now) {
			if task.Repeat == "" {
				taskDate = now // Устанавливаем на текущую дату
			} else {
				taskDateStr, err := logic.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
					return
				}
				taskDate, _ = time.Parse("20060102", taskDateStr)
			}
		}

		// Передаем задачу в бд
		query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
		res, err := db.Exec(query, taskDate.Format("20060102"), task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при добавлении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Получаем ID вставленной записи
		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Ошибка при получении идентификатора задачи"}`, http.StatusInternalServerError)
			return
		}

		// Отправляем успешный ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := fmt.Sprintf(`{"id":"%d"}`, id)
		fmt.Fprintln(w, response)
	}
}

// GetTaskHandler обрабатывает GET-запросы для получения задачи по id.
func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"только метод GET поддерживается"}`, http.StatusMethodNotAllowed)
			return
		}

		// Получаем параметр id из строки запроса
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		// Подготавливаем SQL-запрос
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
		var task Task
		err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при выполнении запроса: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Формируем JSON-ответ
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response, err := json.Marshal(task)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при кодировании ответа"}`, http.StatusInternalServerError)
			return
		}

		// Отправляем JSON-ответ
		w.Write(response)
	}
}

// UpdateTaskHandler обрабатывает PUT-запрос для обновления задачи
func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		// Декодируем JSON-запрос
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		// Проверяем обязательные поля
		if task.ID == "" {
			http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}
		if task.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		// Проверяем и валидируем поле date
		var taskDate time.Time
		if task.Date != "" {
			taskDate, err = time.Parse("20060102", task.Date)
			if err != nil {
				http.Error(w, `{"error":"Дата указана в неверном формате"}`, http.StatusBadRequest)
				return
			}
		} else {
			taskDate = time.Now().Truncate(24 * time.Hour)
		}

		// Проверяем валидность repeat
		if task.Repeat != "" {
			_, err := logic.NextDate(taskDate, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
				return
			}
		}

		// Обновляем задачу в бд
		query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
		res, err := db.Exec(query, taskDate.Format("20060102"), task.Title, task.Comment, task.Repeat, task.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при обновлении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Проверяем, обновилась ли хотя бы одна запись
		rowsAffected, err := res.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		// Отправляем пустой JSON в случае успеха
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}

// DeleteTaskHandler обрабатывает DELETE-запрос для удаления задачи
func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID задачи из параметров запроса
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		// Удаляем задачу из бд
		deleteQuery := "DELETE FROM scheduler WHERE id = ?"
		res, err := db.Exec(deleteQuery, taskID)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при удалении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Проверяем, была ли удалена задача
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при проверке результата удаления: %v"}`, err), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		// Возвращаем пустой JSON в случае успеха
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}

// TaskHandler обрабатывает GET, POST, PUT и DELETE запросы для /api/task
func TaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetTaskHandler(db)(w, r) // Обработка GET-запросов
		case http.MethodPost:
			AddTaskHandler(db)(w, r) // Обработка POST-запросов
		case http.MethodPut:
			UpdateTaskHandler(db)(w, r) // Обработка PUT-запросов
		case http.MethodDelete:
			DeleteTaskHandler(db)(w, r) // Обработка DELETE-запросов
		default:
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	}
}
