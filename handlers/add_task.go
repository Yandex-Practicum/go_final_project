package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/helpers"
	"github.com/FunnyFoXD/go_final_project/models"
)

type insertID struct {
	ID int `json:"id"`
}

// AddTaskHandler is a handler for "/api/task" endpoint.
// It adds a new task to the database.
// It expects a JSON object with the following fields:
// - date: a string in the format "YYYYMMDD"
// - title: a string
// - comment: a string (optional)
// - repeat: a string (optional)
// It returns a JSON object with the following fields:
// - id: an integer representing the ID of the new task
// It returns the following HTTP status codes:
// - 201 Created: the task was successfully added
// - 400 Bad Request: the request body is invalid
// - 500 Internal Server Error: an error occurred while adding the task
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer	
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, `{"error":"can't read body"}`, http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, `{"error":"can't unmarshal body"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"title is empty"}`, http.StatusBadRequest)
		return
	}

	var now = time.Now().Truncate(24 * time.Hour)
	if task.Date == "" || task.Date == "today" || task.Date == "Today" {
		task.Date = now.Format("20060102")
	}

	taskParse, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error":"invalid date format"}`, http.StatusBadRequest)
		return
	}

	if taskParse.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			nextDate, err := helpers.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
				return
			}

			task.Date = nextDate
		}
	}

	var id insertID
	id.ID, err = databases.InsertTask(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Error while inserting task: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
