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
	r.With(api.Auth).Post("/api/task", api.AddTaskHandle)
	r.With(api.Auth).Get("/api/tasks", api.GetTasksHandler)
	r.With(api.Auth).Get("/api/task", api.GetTaskHandler)
	r.With(api.Auth).Put("/api/task", api.PutTaskHandler)
	r.With(api.Auth).Post("/api/task/done", api.TaskDoneHandler)
	r.With(api.Auth).Delete("/api/task", api.DeleteTaskHandler)
	r.Post("/api/signin", api.SignInHandler)

	if err := server.Run(r); err != nil {
		log.Printf("Could not start the server %v\n", err)
	}

}
