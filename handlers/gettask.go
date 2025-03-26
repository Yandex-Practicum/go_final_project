// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go_final_project/database"
	"go_final_project/tasks"
)

// GetTaskByID takes request and get task by ID
func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
	if r.Method != http.MethodGet {
		SendErrorResponse(w, "GetTaskByIDHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get task ID
	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "GetTaskByIDHandler: No ID provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idTask)
	if err != nil {
		SendErrorResponse(w, "GetTaskByIDHandler: Invalid ID format", http.StatusBadRequest)
		return
	}

	var task tasks.Task

	// get task by ID
	task, err = db.GetTaskByID(id)
	if err == sql.ErrNoRows {
		SendErrorResponse(w, fmt.Errorf("GetTaskByIDHandler: failed to find task: %w", err).Error(), http.StatusNotFound)
		return
	} else if err != nil {
		SendErrorResponse(w, fmt.Errorf("GetTaskByIDHandler: failed to retrieve task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	// get JSON response
	response, err := json.Marshal(task)
	if err != nil {
		SendErrorResponse(w, "GetTaskByIDHandler: response JSON creation eror", http.StatusInternalServerError)
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("GetTaskByIDHandler: failed to write response: %v", err)
	}
}
