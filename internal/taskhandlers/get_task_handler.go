package taskhandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go_final_project/internal/domain/entities"
)

// GetTaskHandler - handler for retrieving a task by ID
func GetTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error":"Task ID is missing"}`, http.StatusBadRequest)
			return
		}

		var task entities.Task
		err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
			Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("SQL Error: %v", err)
			http.Error(w, `{"error":"Error retrieving task"}`, http.StatusInternalServerError)
			return
		}

		// Convert task ID to string
		taskIDStr := strconv.FormatInt(task.ID, 10)

		// Send JSON response
		json.NewEncoder(w).Encode(map[string]string{
			"id":      taskIDStr,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		})
	}
}
