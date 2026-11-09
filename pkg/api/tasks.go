package api

import (
	"database/sql"
	"errors"
	"github.com/Yandex-Practicum/go_final_project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(database, 50)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}

func getTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	writeJson(w, task, http.StatusOK)
}
