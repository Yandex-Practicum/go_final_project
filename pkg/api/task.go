package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Stanislav-id05/go_final_project/pkg/db"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writJson(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	writJson(w, http.StatusOK, task)
}

// Функция для записи JSON в ответ
func writJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writJson(w, http.StatusBadRequest, map[string]string{"error": "Некорректные данные"})
		return
	}

	if task.ID == "" {
		writJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if task.Date == "" || !isValidDate(task.Date) {
		writJson(w, http.StatusBadRequest, map[string]string{"error": "Некорректная дата"})
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	if parsedDate, err := time.Parse("20060102", task.Date); err != nil || parsedDate.Before(now) {
		writJson(w, http.StatusBadRequest, map[string]string{"error": "Дата не может быть в прошлом"})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writJson(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	writJson(w, http.StatusOK, map[string]interface{}{})
}

func isValidDate(dateStr string) bool {
	_, err := time.Parse("20060102", dateStr)
	return err == nil
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		http.Error(w, `{"error":"failed to delete task"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
