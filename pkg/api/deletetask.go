package api

import (
	"database/sql"
	"net/http"

	"github.com/Yandex-Practicum/go_final_project/pkg/db"
)

func deleteTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}
	err := db.DeleteTask(database, id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJson(w, map[string]string{})
}
