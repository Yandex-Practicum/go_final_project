package main

import (
	"TODOGo/api"
	"TODOGo/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
)

var (
	portLocal = "7540"
	vebDir    = "./web"
)

func main() {
	config.LoadEnv()
	config.MakeDB()
	defer config.CloseDB()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = portLocal
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Обработка статических файлов
	r.Handle("/*", http.FileServer(http.Dir(vebDir)))

	// Обработка маршрута для входа
	r.Post("/signin", api.SignUpHandler)

	// Группировка маршрутов API
	r.Route("/api", func(r chi.Router) {
		r.Use(api.AuthMiddleware)

		r.Get("/nextdate", api.NextDateHandler)
		r.Post("/task", api.AddTaskHandler)
		r.Get("/tasks", api.GetAllTasksHandler)
		r.Get("/task", api.GetTaskHandler)
		r.Put("/task", api.PutTaskHandler)
		r.Post("/task/done", api.DoneTaskHandler)
		r.Delete("/task", api.DeleteTaskHandler)
	})

	log.Printf("Сервер запущен на порту: %s", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Ошибка загрузки сервера: %s", err)
	}
}
