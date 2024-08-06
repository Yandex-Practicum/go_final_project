// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
)

// DeleteTaskHandler takes request and delete task by ID
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		SendErrorResponse(w, "DeleteTaskHandler: Method not allowed", http.StatusBadRequest)
		return
	}

	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "DeleteTaskHandler: No ID provided", http.StatusBadRequest)
		return
	}

	idTaskParsed, err := strconv.Atoi(idTask)
	if err != nil {
		SendErrorResponse(w, "DeleteTaskHandler: Invalid ID format", http.StatusInternalServerError)
		return
	}
	// delete task by ID
	query := "DELETE FROM scheduler WHERE id = ?"
	result, err := db.Exec(query, idTaskParsed)
	if err != nil {
		SendErrorResponse(w, "DeleteTaskHandler: Error deleting task", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		SendErrorResponse(w, "DeleteTaskHandler: Unable to determine the number of affected rows", http.StatusInternalServerError)
		return
	} else if rowsAffected == 0 {
		SendErrorResponse(w, "DeleteTaskHandler: Task not found", http.StatusInternalServerError)
		return
	}
	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
