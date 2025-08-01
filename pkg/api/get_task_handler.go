package api

import (
	"net/http"

	"github.com/ZhakiyaK/final_project/pkg/db"
)

// Получение задачи
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJsonError(w, "ID is required", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusNotFound)
		return

	}

	writeJson(w, task)
}
