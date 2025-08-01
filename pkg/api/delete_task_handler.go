package api

import (
	"net/http"

	"github.com/ZhakiyaK/final_project/pkg/db"
)

// Удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJsonError(w, "No task with this ID", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{})
}
