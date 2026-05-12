package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"go_final_project/internal/handler"
	"go_final_project/internal/service"
	"go_final_project/internal/store"
	"go_final_project/pkg/api"
)

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
)

func Init(port, webDir string, conn *sql.DB, logs *log.Logger) {
	taskStore := store.NewTaskStore(conn, logs)
	taskService := service.NewTaskService(taskStore, logs)
	taskHandler := handler.NewTaskHandler(taskService, logs)

	srv := http.NewServeMux()
	srv.Handle("/", http.FileServer(http.Dir(webDir)))
	srv.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.NextDate().ServeHTTP(w, r)
		default:
			api.Failed(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})
	srv.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.List().ServeHTTP(w, r)
		default:
			api.Failed(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})
	srv.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.Get().ServeHTTP(w, r)
		case http.MethodPost:
			taskHandler.Add().ServeHTTP(w, r)
		case http.MethodPut:
			taskHandler.Edit().ServeHTTP(w, r)
		case http.MethodDelete:
			taskHandler.Delete().ServeHTTP(w, r)
		default:
			api.Failed(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})
	srv.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskHandler.Done().ServeHTTP(w, r)
		default:
			api.Failed(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":"+port, srv)
}
