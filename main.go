package main

import (
	"fmt"
	"github.com/Ikamenev/database"
	"github.com/Ikamenev/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	database.InitializationDatabase()

	r := chi.NewRouter()
	r.Mount("/", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", handlers.NextDate)
	r.Post("/api/task", handlers.TaskAddPOST)
	r.Get("/api/tasks", handlers.TasksReadGET)
	r.Get("/api/task", handlers.TaskReadGET)
	r.Put("/api/task", handlers.TaskUpdatePUT)
	r.Post("/api/task/done", handlers.TaskDonePOST)
	r.Delete("/api/task", handlers.TaskDELETE)
	fmt.Println("Сервер запущен")

	err := http.ListenAndServe(":7540", r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")

}
