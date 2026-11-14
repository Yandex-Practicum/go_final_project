package api

import (
	"database/sql"
	"net/http"

	"go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}

		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, task)
}
