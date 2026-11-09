package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Yandex-Practicum/go_final_project/pkg/db"
)

func editTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{
			"error": "Не указан заголовок задачи",
		}, http.StatusBadRequest)
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	err = db.UpdateTask(database, &task)
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
