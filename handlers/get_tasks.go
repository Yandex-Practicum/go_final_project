package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/models"
)

type response struct {
	Tasks []models.TaskFromDB `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := databases.GetTasks()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
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
