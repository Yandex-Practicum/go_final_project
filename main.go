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
	"todo_restapi/middlewares"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("init: no .env file found: %v\n", err)
	}
}

func main() {

	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		fmt.Println("no port in .env, will use default port (:7540)")
		port = ":7540"
	}

	storagePath, exists := os.LookupEnv("TODO_DBFILE")
	if !exists {
		fmt.Println("no path in .env, will use default path (./scheduler.db)")
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
	router.Post("/api/signin", taskHandler.Authentication)

	router.With(middlewares.Auth).Route("/api", func(router chi.Router) {
		router.HandleFunc("/task", taskHandler.CRUDTask)
		router.Get("/tasks", taskHandler.GetTasks)
		router.HandleFunc("/task/done", taskHandler.TaskIsDone)
	})

	fmt.Printf("Server is running on port%s...\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("server run error: %v\n", err)
	}
}
