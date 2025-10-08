package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_final_project-main/internal/model"
	"go_final_project-main/internal/nextdate"

	"github.com/jmoiron/sqlx"
)

func AddTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req model.Task
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		// Проверка обязательного поля Title
		if req.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		// Установка текущей даты, если дата не указана
		if req.Date == "" {
			req.Date = time.Now().Format(nextdate.DateFormat)
		}

		// Проверка формата даты
		if _, err := time.Parse(nextdate.DateFormat, req.Date); err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}

		now := time.Now().Format(nextdate.DateFormat)
		if req.Date < now {
			if req.Repeat == "" {
				req.Date = now
			} else {
				nextDate, err := nextdate.NextDate(time.Now(), req.Date, req.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
					return
				}
				req.Date = nextDate
			}
		}

		// Добавление задачи в базу данных
		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.Exec(query, req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при добавлении задачи"}`, http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Ошибка получения ID задачи"}`, http.StatusInternalServerError)
			return
		}

		// Формирование успешного ответа
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response := map[string]interface{}{"id": strconv.FormatInt(id, 10)}
		json.NewEncoder(w).Encode(response)
	}
}

func EditTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			ID      string `json:"id"`
			Date    string `json:"date"`
			Title   string `json:"title"`
			Comment string `json:"comment,omitempty"`
			Repeat  string `json:"repeat,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		// Проверка обязательного поля ID и Title
		if req.ID == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}
		if req.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		// Проверка формата даты
		if _, err := time.Parse(nextdate.DateFormat, req.Date); err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли задача
		id, err := strconv.ParseInt(req.ID, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"Некорректный формат id"}`, http.StatusBadRequest)
			return
		}

		var existingTask model.Task
		err = db.Get(&existingTask, `SELECT * FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		now := time.Now().Format(nextdate.DateFormat)
		if req.Date < now {
			if req.Repeat == "" {
				req.Date = now
			} else {
				nextDate, err := nextdate.NextDate(time.Now(), req.Date, req.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
					return
				}
				req.Date = nextDate
			}
		}

		// Обновление задачи в базе данных
		query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
		result, err := db.Exec(query, req.Date, req.Title, req.Comment, req.Repeat, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		// Возврат пустого JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func GetTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		dbReqId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var task model.Task
		err = db.Get(&task, `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, dbReqId)
		if err != nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}

		response := model.Task{
			ID:      task.ID,
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(response)
	}
}

func AddTaskHandlerWrapper(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AddTaskHandler(db)(w, r)
	}
}
