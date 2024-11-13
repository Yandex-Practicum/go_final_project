package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"pwd/internal/db"
	"pwd/internal/repository"
	"pwd/worker"

	"github.com/go-chi/chi/v5"
)

func main() {
	webDir := "./web"

	db := db.New()
	rep := repository.New(db)
	migration(rep)

	r := chi.NewRouter()

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", worker.NextDateHandler)
	r.Post("/api/task", worker.TaskHandler)
	err := http.ListenAndServe(":7540", nil)
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
