package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go_final_project/pkg/db"
)

func doneTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{
			"error": "Не указан идентификатор",
		}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(database, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, map[string]string{
				"error": "Задача не найдена",
			}, http.StatusNotFound)
			return
		}

		writeInternalError(w, err)
		return
	}

	if task.Repeat == "" {
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
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	err = db.UpdateDate(database, next, id)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
