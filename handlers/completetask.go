// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go_final_project/database"
	"go_final_project/dates"
	"go_final_project/tasks"
)

// DoneTaskHandler takes request and mark task as done by ID
func DoneTaskHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "DoneTaskHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "DoneTaskHandler: No ID provided", http.StatusBadRequest)
		return
	}

	idTaskParsed, err := strconv.Atoi(idTask)
	if err != nil {
		SendErrorResponse(w, "DoneTaskHandler: Invalid ID format", http.StatusBadRequest)
		return
	}

	var task tasks.Task
	task, err = db.GetTaskByID(idTaskParsed)
	if err == sql.ErrNoRows {
		SendErrorResponse(w, fmt.Errorf("DoneTaskHandler: failed to find task: %w", err).Error(), http.StatusNotFound)
		return
	} else if err != nil {
		SendErrorResponse(w, fmt.Errorf("DoneTaskHandler: failed to retrieve task: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now()
	// get task repetition date
	if task.Repeat != "" {
		newTaskDate, err := dates.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			SendErrorResponse(w, "DoneTaskHandler: Invalid repeat pattern", http.StatusBadRequest)
			return
		}

		task.Date = newTaskDate
		// update task if repeat rule set

		err = db.EditTask(task)
		if err != nil {
			SendErrorResponse(w, fmt.Errorf("DoneTaskHandler: failed to update task: %w", err).Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// delete task if repeat rule not set
		err := db.DeleteTask(idTaskParsed)
		if err != nil {
			SendErrorResponse(w, fmt.Errorf("DoneTaskHandler: failed to delete task: %w", err).Error(), http.StatusInternalServerError)
			return
		}
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{}`))
	if err != nil {
		log.Printf("DoneTaskHandler: failed to write response: %v", err)
	}
}
