package main

import (
	"fmt"
	"main/db"
	"main/server"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	database, err := db.CreateDB()
	if err != nil {
		fmt.Println(err)
	}
	defer database.Close()

	webDir := "./web"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
	})
	port := server.GetPort()

	http.HandleFunc("/api/nextdate", server.HandleNextDate)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			server.AddTaskHandler(database)(w, r)
		case http.MethodGet:
			server.GetTaskHandler(database)(w, r)
		case http.MethodPut:
			server.UpdateTaskHandler(database)(w, r)
		case http.MethodDelete:
			server.DeleteTaskHandler(database)(w, r)
		default:
			http.Error(w, "Неверный метод запроса", http.StatusBadRequest)
		}
	})
	http.HandleFunc("/api/tasks", server.GetTasksHandler(database))
	http.HandleFunc("/api/task/done", server.MarkTasksAsDoneHandler(database))

	fmt.Printf("Запуск сервера на порту %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
