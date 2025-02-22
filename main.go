package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"go_final/app/database"
	"go_final/app/handlers"
)

func main() {
	// init database first
	database.InitDB()
	defer database.CloseDB()

	router := http.NewServeMux()

	// register api handlers
	router.HandleFunc("/api/nextdate", handlers.HandlerNewDate)
	router.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		switch r.Method {
		case http.MethodPost:
			handlers.PostTask(w, r)
		case http.MethodGet:
			handlers.GetTask(w, r)
		case http.MethodPut:
			handlers.PutTask(w, r)
		case http.MethodDelete:
			handlers.DeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc("/api/tasks", handlers.GetTasks)
	router.HandleFunc("/api/task/done", handlers.DoneTask)

	webDir := "./web"
	absPath, err := filepath.Abs(webDir)
	if err != nil {
		log.Fatal(err)
	}
	router.Handle("/", http.FileServer(http.Dir(absPath)))

	// start server
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
		log.Printf("Environment variable TODO_PORT not set, using default %s", port)
	}

	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
