package taskhandlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

// DeleteTaskHandler - handler for deleting a task (DELETE /api/task?id=<id>)
func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error":"Task ID is required"}`, http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			log.Printf("Error deleting task: %v", err)
			http.Error(w, `{"error":"Error deleting task"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		}

		log.Printf("Task with ID %d has been deleted", id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}
