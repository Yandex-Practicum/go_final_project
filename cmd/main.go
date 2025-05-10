package main

import (
	"finalproject/pkg/api"
	"finalproject/pkg/db"
	"finalproject/pkg/server"
	"log"

	"github.com/go-chi/chi/v5"
)

func main() {
	_, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Get("/api/nextdate", api.Handler_NextDate)
	r.Post("/api/task", api.AddTaskHandle)
	r.Get("/api/tasks", api.GetTasksHandler)
	r.Get("/api/task", api.GetTaskHandler)
	r.Put("/api/task", api.PutTaskHandler)

	if err := server.Run(r); err != nil {
		log.Printf("Could not start the server %v\n", err)
	}

}
