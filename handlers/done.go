package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/models"
	"go_final_project/utils"
	"net/http"
	"strconv"
	"time"
)

func DoneTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.FormValue("id")
		if idStr == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		var task models.Task
		err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		now := time.Now()
		nextDate, err := utils.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if nextDate == "" {
			_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		} else {
			_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{})
	}
}
