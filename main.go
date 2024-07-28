package main

import (
	"go_final_project/internal/server"
	"go_final_project/internal/server/handler"
	"go_final_project/internal/storage"
	"log"

	"github.com/go-chi/chi"
)

func main() {

	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Handle("/*", handler.GetFront())
	r.Get("/api/nextdate", handler.NextDate)
	r.Post("/api/task", handler.AddTask(db))
	r.Get("/api/tasks", handler.GetTasks(db))
	r.Get("/api/task", handler.GetTask(db))
	r.Put("/api/task", handler.ChangeTask(db))
	r.Post("/api/task/done", handler.DoneTask(db))
	r.Delete("/api/task", handler.DeleteTask(db))

	server := new(server.Server)

	err = server.Run(r)

	if err != nil {
		log.Fatalf("Сервер не запущен %v", err)
		return
	}
}
