package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/models"
)

type response struct {
	Tasks []models.TaskFromDB `json:"tasks"`
}

// GetTasksHandler is a handler for "/api/tasks" endpoint.
// It returns a JSON object with a "tasks" field which is a slice of TaskFromDB structs.
// The returned slice of tasks is limited to 20 tasks.
// If the tasks cannot be retrieved, the function returns an error with the
// following format: "can't get tasks: <error message>".
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := databases.GetTasks(search)
	if err != nil {
		log.Printf("Error while getting tasks: %v", err)
		http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
		return
	}

	response := response{
		Tasks: tasks,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
}
