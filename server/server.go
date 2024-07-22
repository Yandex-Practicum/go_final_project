package server

import (
	"net/http"

	"github.com/Arukasnobnes/go_final_project/handlers"
)

func InitHandlers(h *handlers.Handler) {
	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/nextdate", h.NextDateHandler)
	http.HandleFunc("/api/task", h.TaskHandler)
	http.HandleFunc("/api/task/done", h.TaskDoneHandler)
	http.HandleFunc("/api/tasks", h.TasksListHandler)
}
