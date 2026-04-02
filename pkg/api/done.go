package api

import (
	"net/http"
	"time"

	"todo-list/pkg/db"
)

// DoneTaskHandler обрабатывает POST запросы для /api/task/done
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSONResponse(w, map[string]string{}, http.StatusOK)
		return
	}

	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask := &db.Task{
		ID:      task.ID,
		Date:    nextDate,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	if err := db.UpdateTask(updatedTask); err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]string{}, http.StatusOK)
}
