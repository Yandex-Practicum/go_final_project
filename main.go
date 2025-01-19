package main

import (
	"go_final_project/db"
	"go_final_project/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
		dbPath = filepath.Join(workingDir, "scheduler.db")
	}

	if err := db.SetupDatabase(dbPath); err != nil {
		log.Fatalf("Error with database: %v", err)
	}

	dbConn, err := db.OpenDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	handler := handlers.NewHandler(dbConn)

	http.HandleFunc("/api/tasks", handler.HandleTaskList)
	http.HandleFunc("/api/task/add", handler.HandleAddTask)
	http.HandleFunc("/api/task/get", handler.HandleGetTask)
	http.HandleFunc("/api/task/update", handler.HandleUpdateTask)
	http.HandleFunc("/api/task/delete", handler.HandleDeleteTask)
	http.HandleFunc("/api/task/markdone", handler.HandleMarkTaskDone)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
