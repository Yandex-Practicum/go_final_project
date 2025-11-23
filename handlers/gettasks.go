// package handlers provides API handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go_final_project/database"
	"go_final_project/tasks"
)

type TasksResponse struct {
	Tasks []tasks.Task `json:"tasks"`
}

// GetTasksHandler takes request and provide existing tasks from db
func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *database.Database) {
	if r.Method != http.MethodGet {
		SendErrorResponse(w, "GetTasksHandler: method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var taskList []tasks.Task
	searchStr := r.FormValue("search")

	// get tasks with search query provided
	taskList, err := db.GetTasks(searchStr)
	if err != nil {
		SendErrorResponse(w, fmt.Errorf("GetTasksHandler: failed to get tasks: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	// provide empty array if task list is empty
	if len(taskList) == 0 {
		taskList = []tasks.Task{}
	}

	// get JSON response
	responseMap := map[string][]tasks.Task{"tasks": taskList}
	response, err := json.Marshal(responseMap)
	if err != nil {
		SendErrorResponse(w, "GetTasksHandler: response JSON creation error", http.StatusInternalServerError)
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("GetTasksHandler: failed to write response: %v", err)
	}
}
