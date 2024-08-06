// package main starts server with handlers
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/database"
	"go_final_project/handlers"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

// init loads enviroment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// fileServer is handle func for frontend files
func fileServer(path string) http.Handler {
	fs := http.FileServer(http.Dir(path))
	return http.StripPrefix("/", fs)
}

func main() {

	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	log.Println("Database setup successfully")

	// get enviroment variable for port
	portStr := os.Getenv("TODO_PORT")
	var port int

	// get port number if env variable is empty
	if portStr == "" {
		port = 7540
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
		}
	}

	// Define path for frontend files
	webDir := "./web"

	// register API handlers
	r := chi.NewRouter()

	r.Handle("/*", fileServer(webDir))
	log.Printf("Loaded frontend files from %s\n", webDir)

	r.Get("/api/nextdate", handlers.GetNextDate)
	r.Post("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTaskHandler(w, r, db)
	})
	r.Get("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTasksHandler(w, r, db)
	})
	r.Get("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTaskByIDHandler(w, r, db)
	})
	r.Put("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.EditTaskHandler(w, r, db)
	})
	r.Post("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		handlers.DoneTaskHandler(w, r, db)
	})
	r.Delete("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteTaskHandler(w, r, db)
	})

	// start server
	address := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on http://localhost%s\n", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
		return
	}

}
