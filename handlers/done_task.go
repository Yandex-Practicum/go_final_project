package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/helpers"
)

// DoneTaskHandler is a handler for "/api/task/done" endpoint.
// It marks a task with the given id as done.
// If the id is not given, it returns an error with HTTP status code 400.
// If the id is invalid, it returns an error with HTTP status code 400.
// If there is a database error, it returns an error with HTTP status code 500.
// If the task is successfully marked as done, it returns an empty JSON object with HTTP status code 200.
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, `{"error":"identifier is empty"}`, http.StatusBadRequest)
		return
	}

	task, err := databases.GetTaskByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err = databases.DeleteTask(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := helpers.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
			return
		}

		task.Date = nextDate
		err = databases.UpdateTaskDateByID(id, task.Date)
		if err != nil {
			log.Printf("Error while updating task: %v", err)
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("{}")); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
}
