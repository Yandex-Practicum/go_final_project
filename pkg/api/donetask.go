package api

import (
	"database/sql"
	"github.com/Yandex-Practicum/go_final_project/pkg/db"
	"net/http"
	"time"
)

func doneTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}
	task, err := db.GetTask(database, id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": "Задача не найдена",
		})
		return
	}
	if task.Repeat == "" {
		err := db.DeleteTask(database, id)
		if err != nil {
			writeJson(w, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJson(w, map[string]string{})
		return
	}
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateDate(database, next, id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, map[string]string{})
}
