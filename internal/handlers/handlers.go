package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"go_final-project/internal/db"
	"go_final-project/internal/task"
	"net/http"
)

func GetTasksHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		tasks, err := db.GetTasks(dbase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func AddTaskHandler(dbase *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
			return
		}
		var newTask task.Task
		err := json.NewDecoder(req.Body).Decode(&newTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := db.AddTask(dbase, &newTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newTask.ID = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newTask)
	}
}
