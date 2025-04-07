package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Stanislav-id05/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализация JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		vriteJson(w, map[string]string{"error": "Invalid request payload"})
		return
	}

	// Проверка поля Title
	if task.Title == "" {
		vriteJson(w, map[string]string{"error": "Title cannot be empty"})
		return
	}

	// Проверка даты
	if err := checkDate(&task); err != nil {
		vriteJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		vriteJson(w, map[string]string{"error": "Failed to add task"})
		return
	}

	// Возврат идентификатора в формате JSON
	vriteJson(w, map[string]any{"id": id})
}

func vriteJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)
	var err error
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}
	return nil
}

func afterNow(now, t time.Time) bool {
	return t.Before(now)
}
