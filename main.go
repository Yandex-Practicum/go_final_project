package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	"final/handlers"
	"final/storage"
	"final/tests"
)

func main() {

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	db, err := storage.Createdatabase()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	handlers := handlers.Handlers{db}

	r := chi.NewRouter()
	r.Delete("/api/task", handlers.DeleteTask())
	r.Post("/api/task/done", handlers.TaskDone())
	r.Get("/api/task", handlers.GetTask())
	r.Put("/api/task", handlers.ChangeTask())
	r.Get("/api/tasks", handlers.ReceiveTasks())
	r.Post("/api/task", handlers.AddTask())
	r.Get("/api/nextdate", handlers.NextDate())

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	log.Printf("Сервер слушает порт %s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
		return
	}
}
