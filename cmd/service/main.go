package main

import (
	"log"
	"net/http"

	"go_final_project/internal/database"
	"go_final_project/internal/handler"
	"go_final_project/internal/repository"

	"github.com/go-chi/chi"
)

func main() {

	db := database.New()

	repo := repository.New(db)

	database.Migration(repo)

	handler := handler.New(repo)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", handler.NextDate)
	r.Post("/api/task", handler.AddTask)
	r.Get("/api/tasks", handler.GetTasks)
	r.Get("/api/task", handler.GetTaskById)
	r.Put("/api/task", handler.UpdateTask)
	r.Post("/api/task/done", handler.TaskDone)
	r.Delete("/api/task", handler.DeleteTask)

	if err := http.ListenAndServe(":7540", r); err != nil {
		log.Fatal(err)
	}

}
