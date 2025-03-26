// package handlers provides API handlers
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go_final_project/database"
)

// DeleteTaskHandler takes request and delete task by ID
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
	if r.Method != http.MethodDelete {
		SendErrorResponse(w, "DeleteTaskHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "DeleteTaskHandler: No ID provided", http.StatusBadRequest)
		return
	}

	idTaskParsed, err := strconv.Atoi(idTask)
	if err != nil {
		SendErrorResponse(w, "DeleteTaskHandler: Invalid ID format", http.StatusBadRequest)
		return
	}
	// delete task by ID
	err = db.DeleteTask(idTaskParsed)
	if err != nil {
		SendErrorResponse(w, fmt.Errorf("DeleteTaskHandler: failed to delete task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		log.Printf("DeleteTaskHandler: failed to write response: %v", err)
	}
}
