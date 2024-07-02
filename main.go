package main

import (
	"fmt"
	"main/db"
	"main/server"
	"net/http"
)

func main() {
	db.CreateDB()
	webDir := "./web"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
	})
	port := server.GetPort()

	http.HandleFunc("/api/nextdate", server.HandleNextDate)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			server.AddTaskHandler(w, r)
		case http.MethodGet:
			server.GetTaskHandler(w, r)
		case http.MethodPut:
			server.UpdateTaskHandler(w, r)
		case http.MethodDelete:
			server.DeleteTaskHandler(w, r)
		default:
			http.Error(w, "Неверный метод запроса", http.StatusBadRequest)
		}
	})
	http.HandleFunc("/api/tasks", server.GetTasksHandler)
	http.HandleFunc("/api/task/done", server.MarkTasksAsDoneHandler)

	fmt.Printf("Запуск сервера на порту %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
