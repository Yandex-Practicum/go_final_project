package main

import (
	"database/sql"
	"go_final_project/internal/middleware"
	"go_final_project/internal/taskhandlers"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		log.Println("TODO_PASSWORD is not set, using default password")
		os.Setenv("TODO_PASSWORD", "your_secret_password") // Устанавливаем значение по умолчанию
	}
	log.Println("Current TODO_PASSWORD:", os.Getenv("TODO_PASSWORD"))

	// Connect to the database
	database, err := db.InitDB(db.GetDBFile())
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	defer database.Close()

	// Set up API routes
	apiRouter := http.NewServeMux()

	// Authentication route (public)
	apiRouter.HandleFunc("/api/signin", taskhandlers.SignInHandler)

	// Task management (protected)
	apiRouter.HandleFunc("/api/task", middleware.AuthMiddleware(taskHandler(database)))

	// Fetch tasks (protected)
	apiRouter.HandleFunc("/api/tasks", middleware.AuthMiddleware(taskhandlers.GetTasksHandler(database)))

	// Mark task as done (protected)
	apiRouter.HandleFunc("/api/task/done", middleware.AuthMiddleware(taskhandlers.DoneTaskHandler(database)))

	// Date calculation (public)
	apiRouter.HandleFunc("/api/nextdate", taskhandlers.NextDateHandler)

	// Attach API routes
	http.Handle("/api/", apiRouter)

	// Serve static files
	webDir := getWebDir()
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Directory %s not found!", webDir)
	}

	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	log.Println("Current TODO_PASSWORD:", os.Getenv("TODO_PASSWORD"))

	// Start the server
	port := getPort()
	log.Printf("🚀 Server is running at http://localhost%s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server startup error: %v", err)
	}
}

// taskHandler processes different HTTP methods for /api/task
func taskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskhandlers.AddTaskHandler(db)(w, r)
		case http.MethodGet:
			taskhandlers.GetTaskHandler(db)(w, r)
		case http.MethodPut:
			taskhandlers.EditTaskHandler(db)(w, r)
		case http.MethodDelete:
			taskhandlers.DeleteTaskHandler(db)(w, r)
		default:
			http.Error(w, `{"error":"Method not supported"}`, http.StatusMethodNotAllowed)
		}
	}
}

// getPort returns the server port
func getPort() string {
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		return ":" + envPort
	}
	return ":7540"
}

// getWebDir returns the path to `go_final_project/web`
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

	// Traverse up to find the project root (`go_final_project`)
	for !isProjectRoot(baseDir) {
		baseDir = filepath.Dir(baseDir)
		if baseDir == "/" {
			log.Fatal("Failed to locate the project root `go_final_project`")
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
