package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go_final/app/database"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("id")
	if idStr == "" {
		idStr = r.URL.Query().Get("id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Incorrect id"}`, http.StatusBadRequest)
		return
	}
	// получаем задачу
	task, err := database.GetTaskByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusNotFound)
		return
	}

	response := map[string]string{
		"id":      task.ID,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Error generating JSON"}`, http.StatusInternalServerError)
		return
	}
}
