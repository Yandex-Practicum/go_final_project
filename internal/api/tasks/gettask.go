package tasks

import (
	"net/http"

	"database/sql"
	"final_project/internal/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func GetTasksHandler(res http.ResponseWriter, req *http.Request) {

	search := req.FormValue("search")

	tasks, err := db.GetTasks(50, search)
	if err != nil {
		writeJSON(res, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(res, http.StatusOK, TasksResp{Tasks: tasks})
}
func GetTaskHandler(res http.ResponseWriter, req *http.Request) {
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
	writeJSON(res, http.StatusOK, task)
}
