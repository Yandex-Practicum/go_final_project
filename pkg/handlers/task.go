package handlers

import (
	"database/sql"
	"encoding/json"
	"go-final/pkg/scheduler"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type DBTask struct {
	ID      int    `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

type JSONTask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	ID int64 `json:"id"`
}

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}

		var req TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if req.Title == "" {
			respondWithError(w, http.StatusBadRequest, "Title is required")
			return
		}

		// Устанавливаем now как начало текущего дня (без времени)
		now := time.Now().UTC().Truncate(24 * time.Hour)
		var finalDate time.Time

		// Парсим дату или используем today
		if req.Date == "" || req.Date == "today" || req.Date == now.Format("20060102") {
			finalDate = now
		} else {
			parsedDate, err := time.ParseInLocation("20060102", req.Date, time.UTC)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid date format")
				return
			}
			parsedDate = parsedDate.Truncate(24 * time.Hour) // Обрезаем время
			finalDate = parsedDate

			// Коррекция только для дат в прошлом (сравниваем как даты без времени)
			if finalDate.Before(now) {
				if req.Repeat == "" {
					finalDate = now
				} else {
					next, err := scheduler.NextDate(now, finalDate.Format("20060102"), req.Repeat)
					if err != nil {
						respondWithError(w, http.StatusBadRequest, err.Error())
						return
					}
					finalDate, _ = time.ParseInLocation("20060102", next, time.UTC)
				}
			}
		}

		// Валидация правила повтора (только если дата не today/current)
		if req.Repeat != "" && req.Date != "today" && req.Date != now.Format("20060102") {
			if _, err := scheduler.NextDate(now, finalDate.Format("20060102"), req.Repeat); err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		// Вставка в БД
		res, err := db.Exec(
			`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
			finalDate.Format("20060102"),
			req.Title,
			req.Comment,
			req.Repeat,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		id, _ := res.LastInsertId()
		respondWithSuccess(w, http.StatusCreated, id)
	}
}

// Вспомогательные функции для ответов
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func respondWithSuccess(w http.ResponseWriter, code int, id int64) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(SuccessResponse{ID: id})
}

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		search := r.URL.Query().Get("search")
		var tasks []DBTask
		var err error

		query := "SELECT id, date, title, comment, repeat FROM scheduler"
		args := []interface{}{}
		whereAdded := false
		limit := 50

		if search != "" {
			// Пытаемся распарсить как дату DD.MM.YYYY
			if date, err := time.Parse("02.01.2006", search); err == nil {
				query += " WHERE date = ?"
				args = append(args, date.Format("20060102"))
				whereAdded = true
			} else {
				// Экранируем специальные символы для LIKE
				search = strings.ReplaceAll(search, "%", "\\%")
				search = strings.ReplaceAll(search, "_", "\\_")
				searchTerm := "%" + search + "%"

				if whereAdded {
					query += " AND (title LIKE ? OR comment LIKE ?)"
				} else {
					query += " WHERE (title LIKE ? OR comment LIKE ?)"
				}
				args = append(args, searchTerm, searchTerm)
			}
		}

		query += " ORDER BY date LIMIT ?"
		args = append(args, limit)

		// Логирование для отладки
		log.Printf("Executing query: %s\nArgs: %v", query, args)

		rows, err := db.QueryContext(r.Context(), query, args...)
		if err != nil {
			log.Printf("Database error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var task DBTask
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				log.Printf("Row scan error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
				return
			}
			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Rows error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}

		jsonTasks := make([]JSONTask, 0, len(tasks))
		for _, task := range tasks {
			jsonTasks = append(jsonTasks, JSONTask{
				ID:      strconv.Itoa(task.ID),
				Date:    task.Date,
				Title:   task.Title,
				Comment: task.Comment,
				Repeat:  task.Repeat,
			})
		}

		response := struct {
			Tasks []JSONTask `json:"tasks"`
		}{
			Tasks: jsonTasks,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("JSON encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// GET /api/task
func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Не указан идентификатор"})
			return
		}

		var task DBTask
		err := db.QueryRowContext(r.Context(),
			"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Задача не найдена"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
			return
		}

		json.NewEncoder(w).Encode(JSONTask{
			ID:      strconv.Itoa(task.ID),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}
}

// PUT api/task
func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Валидация ID
		if req.ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Не указан идентификатор"})
			return
		}

		// Валидация заголовка
		if req.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Title is required"})
			return
		}

		// Валидация даты
		now := time.Now().UTC()
		var finalDate time.Time
		if req.Date != "" && req.Date != "today" {
			if _, err := time.Parse("20060102", req.Date); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid date format"})
				return
			}
			parsedDate, _ := time.Parse("20060102", req.Date)
			if parsedDate.Before(now) && req.Repeat == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Дата не может быть в прошлом"})
				return
			}
			finalDate = parsedDate
		} else {
			finalDate = now
		}

		// Валидация правила повторения
		if req.Repeat != "" {
			if _, err := scheduler.NextDate(now, finalDate.Format("20060102"), req.Repeat); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
				return
			}
		}

		// Обновление задачи в БД
		res, err := db.ExecContext(r.Context(),
			"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			finalDate.Format("20060102"),
			req.Title,
			req.Comment,
			req.Repeat,
			req.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err := json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			if err != nil {
				return
			}
			return
		}

		rowsAffected, err := res.RowsAffected()
		if rowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Задача не найдена"})
			if err != nil {
				return
			}
			return
		}

		json.NewEncoder(w).Encode(struct{}{})

	}
}

func MarkDoneHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Не указан идентификатор"})
			return
		}

		var task struct {
			Date   string
			Repeat string
		}
		err := db.QueryRowContext(r.Context(),
			"SELECT date, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.Date, &task.Repeat)

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Задача не найдена"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Ошибка сервера"})
			return
		}

		now := time.Now().UTC()
		//var result sql.Result

		if task.Repeat != "" {
			parsedDate, _ := time.Parse("20060102", task.Date)
			next, err := scheduler.NextDate(now, parsedDate.Format("20060102"), task.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
				return
			}

			// Обновляем дату следующего выполнения
			_, err = db.ExecContext(r.Context(),
				"UPDATE scheduler SET date = ? WHERE id = ?",
				next,
				id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Ошибка сервера"})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(struct{}{})
		} else {
			// Удаляем одноразовую задачу
			_, err := db.ExecContext(r.Context(),
				"DELETE FROM scheduler WHERE id = ?", id)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Ошибка сервера"})
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(struct{}{})
		}
	}
}

// Удаление задачи
func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Не указан идентификатор"})
			return
		}

		result, err := db.ExecContext(r.Context(),
			"DELETE FROM scheduler WHERE id = ?", id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Ошибка сервера"})
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Задача не найдена"})
			return
		}

		json.NewEncoder(w).Encode(struct{}{})
	}
}
