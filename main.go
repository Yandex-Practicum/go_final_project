package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"todo_restapi/internal/handlers"
	"todo_restapi/internal/storage"
)

// go test -run ^TestApp$ ./tests
// go test -run ^TestDB$ ./tests
// go test -run ^TestNextDate$ ./tests
// go test -run ^TestAddTask$ ./tests

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("init: no .env file found: %v\n", err)
	}
}

func main() {

	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		port = ":7540"
	}

	storagePath, exists := os.LookupEnv("TODO_DBFILE")
	if !exists {
		storagePath = "./scheduler.db"
	}

	database, err := storage.OpenStorage(storagePath)
	if err != nil {
		log.Fatalf("OpenStorage: %v", err)
	}

	taskHandler := handlers.NewTaskHandler(database)

	router := chi.NewRouter()

	router.Get("/", func(write http.ResponseWriter, request *http.Request) {
		http.ServeFile(write, request, "web/index.html")
	})

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("web"))))

	router.Get("/api/nextdate", handlers.NextDateHandler)
	router.HandleFunc("/api/task", taskHandler.AddTask)

	fmt.Printf("Server is running on port%s...\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("server run error: %v\n", err)
	}
}
