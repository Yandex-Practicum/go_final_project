package api

import (
	"net/http"
	"strconv"

	"go_final_project/pkg/db"
)

// TasksResp — структура JSON-ответа
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Параметры: limit (необязательно), search (необязательно)
	q := r.URL.Query()
	limit := 50
	if l := q.Get("limit"); l != "" {
		if i, err := strconv.Atoi(l); err == nil && i > 0 {
			limit = i
		}
	}
	search := q.Get("search")

	tasks, err := db.Tasks(limit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
