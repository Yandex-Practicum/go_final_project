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

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	return port
}

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

	router.Get("/api/nextdate", handlers.NewDateHandler)
	router.Get("/api/tasks", handlers.GetTasksHandler)
	router.Get("/api/task", handlers.GetTaskHandler)

	router.Post("/api/task", handlers.AddTaskHandler)
	
	router.Put("/api/task", handlers.UpdateTaskHandler)

	log.Printf("starting listen server on port %s", port)
	for err := http.ListenAndServe(":"+port, router); err != nil; {
		log.Fatalf("start server error: %s", err.Error())
	}
}
