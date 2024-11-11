package main

import (
	"log"
	"net/http"

	"pwd/database"
	"pwd/handlers"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {
	webDir := "./web"

	database.ConnectDb()

	r := chi.NewRouter()
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", handlers.NextDateHandler)
	r.Post("/api/task", handlers.TaskHandler)
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal(err)
	}
}
