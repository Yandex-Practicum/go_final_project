package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_final/app/database"
	"go_final/app/tasks"
)

func DoneTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"Task id not specified"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"Incorrect task id"}`, http.StatusBadRequest)
		return
	}

	task, err := database.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Error getting task from database"}`, http.StatusInternalServerError)
		}
		return
	}

	if task.Repeat == "" {
		if err := database.DeleteTaskByID(id); err != nil {
			http.Error(w, `{"error":"Error deleting task"}`, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Error calculating next recurrence date"}`, http.StatusInternalServerError)
			return
		}

		if err := database.UpdateTaskDate(uint64(id), nextDate); err != nil {
			http.Error(w, `{"error":"Error task update"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}
