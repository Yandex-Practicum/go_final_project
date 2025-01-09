package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/helpers"
	"github.com/FunnyFoXD/go_final_project/models"
)

// UpdateTaskHandler is a handler for "/api/task" endpoint.
// It updates a task with the given id in the database.
// It expects a PUT request with a JSON object with the following fields:
// - id: a string
// - date: a string (optional)
// - title: a string
// - comment: a string (optional)
// - repeat: a string (optional)
// It returns a JSON object with the following fields:
// - id: an integer representing the ID of the updated task
// It returns the following HTTP status codes:
// - 200 OK: the task was successfully updated
// - 400 Bad Request: the request body is invalid
// - 404 Not Found: the task with the given id cannot be found
// - 500 Internal Server Error: an error occurred while updating the task
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, `{"error":"can't read body"}`, http.StatusBadRequest)
		return
	}

	var updatedTask models.TaskFromDB
	if err := json.Unmarshal(buf.Bytes(), &updatedTask); err != nil {
		http.Error(w, `{"error":"can't unmarshal body"}`, http.StatusBadRequest)
		return
	}

	if updatedTask.Title == "" {
		http.Error(w, `{"error":"title is empty"}`, http.StatusBadRequest)
		return
	}

	var now = time.Now().Truncate(24 * time.Hour)
	if updatedTask.Date == "" || updatedTask.Date == "today" || updatedTask.Date == "Today" {
		updatedTask.Date = now.Format("20060102")
	}

	taskParse, err := time.Parse("20060102", updatedTask.Date)
	if err != nil {
		http.Error(w, `{"error":"invalid date format"}`, http.StatusBadRequest)
		return
	}

	if taskParse.Before(now) {
		if updatedTask.Repeat == "" {
			updatedTask.Date = now.Format("20060102")
		} else {
			nextDate, err := helpers.NextDate(now, updatedTask.Date, updatedTask.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
				return
			}

			updatedTask.Date = nextDate
		}
	}

	err = databases.UpdateTaskByID(updatedTask)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error while updating task: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("{}")); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
}
