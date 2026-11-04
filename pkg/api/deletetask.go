package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Задача не найдена",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
