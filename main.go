package main

import (
	"go_final_project/server"
	"go_final_project/server/handler"
	"go_final_project/storage"
	"log"

	"github.com/go-chi/chi"
)

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close storage: %v", err)
		}
	}()

	r := chi.NewRouter()

	r.Handle("/*", handler.GetFront())

	r.Get("/api/nextdate", handler.GetNextDate)
	r.Post("/api/task", handler.AddTask(db))
	r.Get("/api/tasks", handler.GetTasks(db))
	r.Get("/api/task", handler.GetTask(db))
	r.Put("/api/task", handler.UpdateTask(db))
	r.Post("/api/task/done", handler.DoneTask(db))
	r.Delete("/api/task", handler.DeleteTask(db))

	server := new(server.Server)
	if err := server.Run(r); err != nil {
		log.Fatalf("Server can't start: %v", err)
		return
	}
}
