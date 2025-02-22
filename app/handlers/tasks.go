package handlers

import (
	"encoding/json"
	"net/http"

	"go_final/app/database"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	const limit = 50

	tasks, err := database.FindInDB(search, limit)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	var tasksForResponse []map[string]string
	for _, task := range tasks {
		tasksForResponse = append(tasksForResponse, map[string]string{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		})
	}

	if tasksForResponse == nil {
		tasksForResponse = []map[string]string{}
	}

	response := map[string]interface{}{
		"tasks": tasksForResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Error generationg JSON"}`, http.StatusInternalServerError)
		return
	}
}
