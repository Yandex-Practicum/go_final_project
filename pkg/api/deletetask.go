package api

import (
	"database/sql"
	"errors"
	"github.com/Yandex-Practicum/go_final_project/pkg/db"
	"net/http"
)

func deleteTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{
			"error": "Не указан идентификатор",
		}, http.StatusBadRequest)
		return
	}
	err := db.DeleteTask(database, id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJson(w, map[string]string{
				"error": "Задача не найдена",
			}, http.StatusNotFound)
			return
		}

		writeInternalError(w, err)
		return
	}
	writeJson(w, map[string]string{}, http.StatusOK)
}
