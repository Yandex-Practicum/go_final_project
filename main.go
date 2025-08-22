package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort = "7540"
	webDir      = "./web"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(".env file not found")
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/tasks", auth(getTasksHandler(db)))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler(db)))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler(db)))

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func taskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskHandler(db)(w, r)
		case http.MethodPost:
			createTaskHandler(db)(w, r)
		case http.MethodPut:
			updateTaskHandler(db)(w, r)
		case http.MethodDelete:
			deleteTaskHandler(db)(w, r)
		default:
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	}
}
