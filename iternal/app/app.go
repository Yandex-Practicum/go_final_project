package app

import (
	"Go/iternal/database"
	handlers "Go/iternal/transport/rest"

	"github.com/go-chi/chi"

	"fmt"
	"net/http"
)

func Run() {
	r := chi.NewRouter()
	_, err := database.CreateDB()
	if err != nil {
		panic(err)
	}
	fmt.Println("Запускаем сервер!")

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.HandleFunc("/api/task", handlers.Task)
	r.HandleFunc("/api/nextdate", handlers.NextDeadLine)
	r.HandleFunc("/api/tasks", handlers.GetTasks)

	err = http.ListenAndServe(":7540", r)
	if err != nil {
		panic(err)
	}
}
