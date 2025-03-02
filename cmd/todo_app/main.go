package main

import (
	"go_final_project/internal/taskhandlers"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Connect to the database
	database, err := db.InitDB(db.GetDBFile())
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	defer database.Close()

	// Set up API routes
	apiRouter := http.NewServeMux()

	// Handler for all methods on /api/task
	apiRouter.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskhandlers.AddTaskHandler(database)(w, r)
		case http.MethodGet:
			taskhandlers.GetTaskHandler(database)(w, r)
		case http.MethodPut:
			taskhandlers.EditTaskHandler(database)(w, r)
		case http.MethodDelete:
			taskhandlers.DeleteTaskHandler(database)(w, r)
		default:
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
		}
	})

	apiRouter.HandleFunc("/api/tasks", taskhandlers.GetTasksHandler(database))
	apiRouter.HandleFunc("/api/nextdate", taskhandlers.NextDateHandler)

	// Added task completion handler
	apiRouter.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			taskhandlers.DoneTaskHandler(database)(w, r)
		} else {
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
		}
	})

	// Attach API routes
	http.Handle("/api/", apiRouter)

	// Serve static files (fixed path)
	webDir := getWebDir()
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Directory %s not found!", webDir)
	}

	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Start the server
	port := getPort()
	log.Printf("🚀 Server is running at http://localhost%s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}

// getPort returns the server port
func getPort() string {
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		return ":" + envPort
	}
	return ":7540"
}

// getWebDir returns the path to `go_final_project_ref/web`
func getWebDir() string {
	if envDir := os.Getenv("WEB_DIR"); envDir != "" {
		log.Printf("📂 Using web client directory from environment variable: %s", envDir)
		return envDir
	}

	// Get the project root directory path
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting the working directory: %v", err)
	}

	// Traverse up to find the project root (go_final_project_ref)
	for !isProjectRoot(baseDir) {
		baseDir = filepath.Dir(baseDir)
		if baseDir == "/" {
			log.Fatal("Failed to locate the project root go_final_project_ref")
		}
	}

	// Construct the `web` path
	webPath := filepath.Join(baseDir, "web")
	log.Printf("📂 Expected web path: %s", webPath)

	return webPath
}

// isProjectRoot checks if the given path is the project root
func isProjectRoot(path string) bool {
	_, err := os.Stat(filepath.Join(path, "go.mod"))
	return err == nil
}
