package tasks

import (
	"database/sql"
	"net/http"
	"time"

	"final_project/internal/api"
	"final_project/internal/db"
)

func TaskDoneHandler(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		writeJSON(res, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(res, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
			return
		}
		writeJSON(res, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(res, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	} else {
		now := time.Now()
		next, err := api.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(res, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		err = db.UpdateDate(id, next)
		if err != nil {
			writeJSON(res, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
	writeJSON(res, http.StatusOK, map[string]string{})
}
