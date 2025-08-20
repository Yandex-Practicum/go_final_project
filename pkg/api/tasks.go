package api

import (
	"Final/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	limit := 50
	search := r.URL.Query().Get("search")
	var tasks []*db.Task
	var err error

	if search == "" {
		tasks, err = db.Tasks(limit)
	} else {
		tasks, err = db.SearchTasks(search, limit)
	}

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
