package route

import (
	"github.com/gorilla/mux"
	"go_final_project/internal/handler"
	"net/http"
)

func SetupRouter() http.Handler {
	// Настройка роутера
	r := mux.NewRouter()
	r.HandleFunc("/api/nextdate", handler.NextDateHandler).Methods("GET")
	r.HandleFunc("/api/task", handler.AddTaskHandler).Methods("POST")
	r.HandleFunc("/api/tasks", handler.GetTasksListHandler).Methods("GET")
	r.HandleFunc("/api/task", handler.GetTaskHandler).Methods("GET")

	// Обработка статических файлов
	staticDir := "./web/"
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	return r
}
