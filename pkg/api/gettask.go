package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

// GET /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}
	// ищем задачу в бд

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	writeJson(w, task, http.StatusOK)

}
