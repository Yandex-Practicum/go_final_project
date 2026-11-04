package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Задача не найдена",
		})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
