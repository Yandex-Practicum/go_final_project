package main

import (
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	webDir := "./web"

	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.HandleFunc("/api/task/done", api.DoneTaskHandler)
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server started on port", port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
