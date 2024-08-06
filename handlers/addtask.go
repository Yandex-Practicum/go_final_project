package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/tasks"
	"net/http"
	"time"
)

// AddTaskHandler
func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "AddTaskHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// decode JSON
	var task tasks.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: JSON deserialization error", http.StatusBadRequest)
		return
	}

	// check task title existence
	if task.Title == "" {
		SendErrorResponse(w, "AddTaskHandler: Task title not specified", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	// check date format
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Invalid date format", http.StatusBadRequest)
		return
	}

	// get task repetition date
	if task.Repeat != "" {
		dateCheck, err := tasks.NextDate(time.Now(), task.Date, task.Repeat)
		if dateCheck == "" && err != nil {
			SendErrorResponse(w, "AddTaskHandler: Invalid repeat rule", http.StatusBadRequest)
			return
		}
	}

	now := time.Now()
	if date.Before(now) {
		if task.Repeat == "" || date.Truncate(24*time.Hour) == date.Truncate(24*time.Hour) {
			task.Date = time.Now().Format("20060102")
		} else {
			dateStr := date.Format("20060102")
			nextDate, err := tasks.NextDate(now, dateStr, task.Repeat)
			if err != nil {
				SendErrorResponse(w, "AddTaskHandler: Invalid repeat rule", http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	// add task
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Error executing db query", http.StatusInternalServerError)
		return
	}

	// get new task ID
	id, err := res.LastInsertId()
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Error getting new task ID", http.StatusInternalServerError)
		return
	}

	task.Id = fmt.Sprint(id)

	taskId := map[string]interface{}{"id": id}
	response, err := json.Marshal(taskId)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: response JSON creation  error", http.StatusInternalServerError)
		return
	}
	// send response
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
