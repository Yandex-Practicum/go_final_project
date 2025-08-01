package api

import (
	"net/http"
	"time"

	"github.com/ZhakiyaK/final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJsonError(w, "ID is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(w, "Task not found", http.StatusBadRequest)
		return
	}

	// Если нет правил повторения, удаляем таску
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJson(w, map[string]string{})
		return
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)

	next, err := NextDate(now, task.Date, task.Repeat)

	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.UpdateDate(id, next); err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{})

}
