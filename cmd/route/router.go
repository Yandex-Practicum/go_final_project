package route

import (
	"go_final_project/internal/handler"
	"go_final_project/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() http.Handler {
	// Настройка роутера
	r := mux.NewRouter()
	r.HandleFunc("/api/signin", handler.Signin).Methods("POST")
	r.HandleFunc("/api/nextdate", handler.NextDateHandler).Methods("GET")

	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.AuthMiddleware)
	api.HandleFunc("/api/task", handler.AddTaskHandler).Methods("POST")
	api.HandleFunc("/api/tasks", handler.GetTasksListHandler).Methods("GET")
	api.HandleFunc("/api/task", handler.GetTaskHandler).Methods("GET")
	api.HandleFunc("/api/task", handler.UpdateTask).Methods("PUT")
	api.HandleFunc("/api/task/done", handler.DoneDeleteTask).Methods("POST")
	api.HandleFunc("/api/task", handler.DoneDeleteTask).Methods("DELETE")

	// Обработка статических файлов
	staticDir := "./web/"
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	return r
}
