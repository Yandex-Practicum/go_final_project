package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		writeError(w, "Метод не поддерживается")
		return
	}

	tasks, err := db.Tasks(1000)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	today := time.Now().Format("20060102")
	for i := range tasks {
		if tasks[i].Date == "" {
			tasks[i].Date = today
		}
	}

	writeJson(w, map[string]any{
		"tasks": tasks,
	})

}

func writeError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
