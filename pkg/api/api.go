package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func Init(database *sql.DB) {
	http.HandleFunc("/api/nextdate", nextDateHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		taskHandler(database, w, r)
	})
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasksHandler(database, w, r)
	})
	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
		doneTaskHandler(database, w, r)
	})
}

func taskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(database, w, r)
	case http.MethodGet:
		getTaskHandler(database, w, r)
	case http.MethodPut:
		editTaskHandler(database, w, r)
	case http.MethodDelete:
		deleteTaskHandler(database, w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
