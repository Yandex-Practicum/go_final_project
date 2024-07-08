package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/db"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(db.DateFormat, nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: Неверный формат даты 'now': %s", err), http.StatusBadRequest)
		return
	}

	nextDate, err := db.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка: %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", nextDate)
}

func AddTaskHandler(s *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task db.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			writeError(w, "Ошибка десериализации JSON: ", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			writeError(w, "Не указан заголовок задачи", http.StatusBadRequest)
			return
		}

		if task.Date == "" {
			task.Date = time.Now().Format(db.DateFormat)
		} else {
			_, err = time.Parse(db.DateFormat, task.Date)
			if err != nil {
				writeError(w, "Неверный формат даты, ожидается формат "+db.DateFormat, http.StatusBadRequest)
				return
			}
		}

		if task.Repeat == "" {
			if task.Date < time.Now().Format(db.DateFormat) {
				task.Date = time.Now().Format(db.DateFormat)
			}
		} else {
			nextDate, err := db.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeError(w, "Неверный формат правила повторения: "+err.Error(), http.StatusBadRequest)
				return
			}
			if task.Date < time.Now().Format(db.DateFormat) {
				task.Date = nextDate
			}
		}

		id, err := db.AddTask(s, task)
		if err != nil {
			writeError(w, "Ошибка добавления задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writeResponse(w, map[string]int64{"id": id}, http.StatusCreated)
	}
}

func GetTaskHandler(s *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeError(w, "Не указан ID задачи", http.StatusBadRequest)
			return
		}
		task, err := db.GetTaskByID(s, id)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeResponse(w, task, http.StatusOK)
	}
}
func UpdateTaskHandler(s *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updatedTask db.Task
		err := json.NewDecoder(r.Body).Decode(&updatedTask)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedTask, err = updatedTask.Check()
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.PutTask(s, updatedTask)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeResponse(w, map[string]string{}, http.StatusCreated)
	}
}
func GetTasksHandler(s *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		tasksList, err := db.GetTasks(s, search)
		if err != nil {
			writeError(w, "Ошибка поиска задач: ", http.StatusInternalServerError)
			return
		}

		writeResponse(w, map[string][]db.Task{"tasks": tasksList}, http.StatusOK)
	}
}
func DeleteTaskHandler(s *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeError(w, "Не указан ID задачи", http.StatusBadRequest)
			return
		}
		err := db.DeleteTask(s, id)
		if err != nil {
			writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
			return
		}

		writeResponse(w, struct{}{}, http.StatusOK)
	}
}

func MarkTasksAsDoneHandler(s *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			writeError(w, "Не указан ID задачи", http.StatusBadRequest)
			return
		}

		task, err := db.GetTaskByID(s, id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, "Задача не найдена", http.StatusNotFound)
			} else {
				writeError(w, "Ошибка получения задачи", http.StatusInternalServerError)
			}
			return
		}

		if task.Repeat != "" {
			nextDate, err := db.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeError(w, "Ошибка расчета следующей даты", http.StatusInternalServerError)
				return
			}
			task.Date = nextDate
			err = db.PutTask(s, task)
			if err != nil {
				writeError(w, "Ошибка обновления задачи", http.StatusInternalServerError)
				return
			}
		} else {
			err = db.DeleteTask(s, id)
			if err != nil {
				writeError(w, "Ошибка удаления задачи", http.StatusInternalServerError)
				return
			}
		}

		writeResponse(w, struct{}{}, http.StatusOK)
	}
}
func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
