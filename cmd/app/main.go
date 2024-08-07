package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/internal/db"
	"go_final_project/internal/handlers"
)

var (
	defaultPort = 7540
	webDir      = "./web"
	dbInstance  *sql.DB
)

func main() {
	port := defaultPort
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	var err error
	dbInstance, err = db.InitDatabase()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer dbInstance.Close()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	http.HandleFunc("/api/signin", handlers.SignHandler())
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/tasks", auth(handlers.GetTaskHandler(dbInstance)))
	http.HandleFunc("/api/task", auth(handlers.TaskHandler(dbInstance)))
	http.HandleFunc("/api/task/done", auth(handlers.TaskDoneHandler(dbInstance)))

	log.Printf("Сервер запущен на порту %d\n", port)
	log.Printf("Обслуживание файлов из каталога: %s\n", webDir)
	listenAddr := fmt.Sprintf("localhost:%d", port)
	log.Printf("Запуск сервера на %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			if !handlers.IsTokenValid(jwt) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
