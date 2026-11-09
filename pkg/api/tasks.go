package api

import (
	"database/sql"
	"net/http"

	"github.com/Yandex-Practicum/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(database, 50)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func getTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJson(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(database, id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJson(w, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}

		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, task)
}
