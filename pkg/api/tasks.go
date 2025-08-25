package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

var Limit = 100

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(Limit)
	if err != nil {
		http.Error(w, "Limit 100 entries", http.StatusInternalServerError)
		return
	}

	today := time.Now().Format(Layout)
	for i := range tasks {
		if tasks[i].Date == "" {
			tasks[i].Date = today
		}
	}

	writeJson(w, map[string]any{
		"tasks": tasks,
	})

}

func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return
	}
}
