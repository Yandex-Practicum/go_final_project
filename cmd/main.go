package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/FunnyFoXD/go_final_project/databases"
	"github.com/FunnyFoXD/go_final_project/handlers"
	"github.com/FunnyFoXD/go_final_project/tests"
)

// getPort returns the port number to be used for the server.
//
// If TODO_PORT environment variable is set, it will be used, otherwise
// tests.Port will be used as a default value.
func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	return port
}

// main is the entry point for the server application.
//
// It changes the current working directory to "./cmd", initializes the 
// router, and sets up the database. The function then defines the routes 
// for handling HTTP requests, including sign-in, task management, and 
// authorization. Finally, it starts the HTTP server on the specified port.

func main() {
	err := os.Chdir("./cmd")
	if err != nil {
		log.Fatalf("can't change directory: %s", err.Error())
	}

	router := chi.NewRouter()
	port := getPort()

	err = databases.CreateDB()
	if err != nil {
		log.Fatalf("can't create database: %s", err.Error())
	}

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("../web"))))

	router.Post("/api/signin", handlers.SigninHandler)
	router.Post("/api/task", handlers.Authorization(handlers.AddTaskHandler))
	router.Post("/api/task/done", handlers.Authorization(handlers.DoneTaskHandler))

	router.Get("/api/nextdate", handlers.NextDateHandler)
	router.Get("/api/tasks", handlers.Authorization(handlers.GetTasksHandler))
	router.Get("/api/task", handlers.Authorization(handlers.GetTaskHandler))

	router.Put("/api/task", handlers.Authorization(handlers.UpdateTaskHandler))

	router.Delete("/api/task", handlers.Authorization(handlers.DeleteTaskHandler))

	log.Printf("starting listen server on port %s", port)
	for err := http.ListenAndServe(":"+port, router); err != nil; {
		log.Fatalf("start server error: %s", err.Error())
	}
}
