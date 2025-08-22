package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "missing id"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "task not found"}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": "delete failed"}, http.StatusInternalServerError)
			return
		}
		writeJson(w, map[string]any{}, http.StatusOK)
		return
	}

	dateParsed, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid task date"}, http.StatusInternalServerError)
		return
	}

	next, err := nextdate.NextDate(dateParsed, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJson(w, map[string]string{"error": "update failed"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
