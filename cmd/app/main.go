package main

import (
	"database/sql"
	"fmt"
	"go_final_project/internal/middleware"
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

	auth := middleware.NewAuth()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	http.HandleFunc("/api/signin", handlers.SignHandler())
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/tasks", auth.Handle(handlers.GetTaskHandler(dbInstance)))
	http.HandleFunc("/api/task", auth.Handle(handlers.TaskHandler(dbInstance)))
	http.HandleFunc("/api/task/done", auth.Handle(handlers.TaskDoneHandler(dbInstance)))

	log.Printf("Сервер запущен на порту %d\n", port)
	log.Printf("Обслуживание файлов из каталога: %s\n", webDir)
	listenAddr := fmt.Sprintf("localhost:%d", port)
	log.Printf("Запуск сервера на %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
