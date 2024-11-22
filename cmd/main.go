package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pwd/pkg/db"
	"pwd/worker"

	"pwd/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	webDir := "./web"

	d, err := db.New()
	if err != nil {
		panic(err)
	}

	rep := repository.New(d)
	migration(rep)

	r := chi.NewRouter()

	tc := worker.NewTaskController(d)

	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", tc.NextDateHandler)
	r.Get("/api/tasks", tc.GetTaskHandler)
	r.Post("/api/task", tc.PostTaskHandler)
	r.Get("/api/task", tc.GetTaskId)
	r.Put("/api/task/{id}", tc.UpdateTaskId)
	err = http.ListenAndServe(":7540", r)
	if err != nil {
		log.Fatal(err)
	}
}

func migration(rep *repository.Repository) {

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		if err := rep.CreateScheduler(); err != nil {
			log.Fatal(err)
		}
	}
}
