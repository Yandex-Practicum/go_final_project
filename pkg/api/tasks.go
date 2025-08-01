package api

import (
	"net/http"

	"github.com/ZhakiyaK/final_project/pkg/db"
)

var Limit = 50

type TasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.Tasks(Limit)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, TasksResponse{
		Tasks: tasks})
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
