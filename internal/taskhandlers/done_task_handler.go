package taskhandlers

import (
	"database/sql"
	"go_final_project/internal/domain/services"
	"log"
	"net/http"
	"time"
)

// DoneTaskHandler - handler for marking a task as completed (POST /api/task/done)
func DoneTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
			return
		}

		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error":"Task ID is required"}`, http.StatusBadRequest)
			return
		}

		// Retrieve the task date and repeat rule
		var originalDateStr, repeatRule string
		err := db.QueryRow("SELECT date, repeat FROM scheduler WHERE id = ?", taskID).Scan(&originalDateStr, &repeatRule)
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Error retrieving task: %v", err)
			http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
			return
		}

		// If repeat is empty → delete the task
		if repeatRule == "" {
			_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", taskID)
			if err != nil {
				log.Printf("Error deleting task: %v", err)
				http.Error(w, `{"error":"Error deleting task"}`, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
			return
		}

		// Calculate the next execution date, passing `time.Now()` and the original date as a string
		nextDate, err := services.NextDate(time.Now(), originalDateStr, repeatRule)
		if err != nil {
			log.Printf("Error calculating the next date: %v", err)
			http.Error(w, `{"error":"Error calculating the next date"}`, http.StatusInternalServerError)
			return
		}

		// Update the task date in the database
		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, taskID)
		if err != nil {
			log.Printf("Error updating task date: %v", err)
			http.Error(w, `{"error":"Error updating task date"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}
