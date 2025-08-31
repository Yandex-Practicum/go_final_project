package api

import (
	"net/http"
	"strconv"

	"go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []map[string]string `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]any{"error": "method not allowed"})
		return
	}

	items, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	out := make([]map[string]string, 0, len(items))
	for _, it := range items {
		out = append(out, map[string]string{
			"id":      strconv.FormatInt(it.ID, 10),
			"date":    it.Date,
			"title":   it.Title,
			"comment": it.Comment,
			"repeat":  it.Repeat,
		})
	}

	writeJSON(w, TasksResp{Tasks: out})
}
