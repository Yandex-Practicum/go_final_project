// package handlers provides API handlers
package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/tasks"
	"net/http"
	"strconv"
	"time"
)

// GetTaskByID takes request and get task by ID
func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		SendErrorResponse(w, "GetTaskByID: Method not allowed", http.StatusBadRequest)
		return
	}

	// get task ID
	idTask := r.FormValue("id")
	if idTask == "" {
		SendErrorResponse(w, "GetTaskByID: No ID provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idTask)
	if err != nil {
		SendErrorResponse(w, "GetTaskByID: Invalid ID format", http.StatusBadRequest)
		return
	}

	var task tasks.Task

	// get task by ID
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		SendErrorResponse(w, "GetTaskByID: Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		SendErrorResponse(w, "GetTaskByID: Error executing db query", http.StatusInternalServerError)
		return
	}

	// get JSON response
	response, err := json.Marshal(task)
	if err != nil {
		SendErrorResponse(w, "GetTaskByID: response JSON creation eror", http.StatusInternalServerError)
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// EditTaskHandler takes request and edit task by ID
func EditTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "EditTaskHandler: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// decode JSON
	var task tasks.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: JSON deserialization error", http.StatusBadRequest)
		return
	}

	// check task ID
	if task.Id == "" {
		SendErrorResponse(w, "EditTaskHandler: Task ID not found", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(task.Id)
	if err != nil || id <= 0 {
		SendErrorResponse(w, "EditTaskHandler: Invalid task ID", http.StatusBadRequest)
		return
	}

	// check task title
	if task.Title == "" {
		SendErrorResponse(w, "EditTaskHandler: Task title must be specified", http.StatusBadRequest)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	// parse task Date
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: Invalid date format", http.StatusBadRequest)
		return
	}

	// check task repetition rule
	if task.Repeat != "" {
		if _, err := strconv.Atoi(task.Repeat[2:]); err != nil || (task.Repeat[0] != 'd' && task.Repeat[0] != 'y') {
			SendErrorResponse(w, "EditTaskHandler: Invalid task repetition format", http.StatusBadRequest)
			return
		}
	}

	// check task existence
	var idTask int
	err = db.QueryRow("SELECT id FROM scheduler WHERE id = ?", task.Id).Scan(&idTask)
	if err == sql.ErrNoRows {
		SendErrorResponse(w, "EditTaskHandler: Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		SendErrorResponse(w, "EditTaskHandler: Error checking task existence", http.StatusInternalServerError)
		return
	}

	// update task
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat =? WHERE id = ?"

	_, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: Task not found", http.StatusInternalServerError)
		return
	}

	// get JSON response
	response, err := json.Marshal(struct{}{})
	if err != nil {
		SendErrorResponse(w, "EditTaskHandler: response JSON creation error", http.StatusInternalServerError)
		return
	}
	// send response
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
