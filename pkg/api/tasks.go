package api

import (
	"final_sprint/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	search := r.FormValue("search")

	const defaultLimit = 50
	tasks, err := db.Tasks(defaultLimit, search)

	if err != nil {
		writeError(w, err.Error())
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
