package api

import (
	"encoding/json"
	"net/http"

	"go_final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	task.Date, err = checkDate(task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}
