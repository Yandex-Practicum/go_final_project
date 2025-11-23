// package handlers provides API handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go_final_project/constants"
	"go_final_project/database"
	"go_final_project/dates"
	"go_final_project/tasks"
)

// AddTaskHandler
func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
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
		task.Date = time.Now().Format(constants.DateFormat)
	}

	// check date format
	date, err := time.Parse(constants.DateFormat, task.Date)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Invalid date format", http.StatusBadRequest)
		return
	}

	// get task repetition date
	if task.Repeat != "" {
		dateCheck, err := dates.NextDate(time.Now(), task.Date, task.Repeat)
		if dateCheck == "" && err != nil {
			SendErrorResponse(w, "AddTaskHandler: Invalid repeat rule", http.StatusBadRequest)
			return
		}
	}

	task.Date, err = dates.GetTaskRepetitionDate(task.Repeat, date)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: Invalid repeat rule", http.StatusBadRequest)
		return
	}

	// add task
	idTask, err := db.AddTask(task)
	if err != nil {
		SendErrorResponse(w, fmt.Errorf("AddTaskHandler: failed to add task: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	task.Id = fmt.Sprint(idTask)

	taskId := map[string]interface{}{"id": task.Id}
	response, err := json.Marshal(taskId)
	if err != nil {
		SendErrorResponse(w, "AddTaskHandler: response JSON creation  error", http.StatusInternalServerError)
		return
	}
	// send response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("AddTaskHandler: failed to write response: %v", err)
	}
}
