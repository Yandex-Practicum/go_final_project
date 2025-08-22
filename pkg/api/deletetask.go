package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJson(w, map[string]string{"error": "Missing task ID"}, http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Task not found"}, http.StatusNotFound)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
