package main

import (
	"database/sql"
	"fmt"
	"go_final_project/internal/handlers/auth"
	"go_final_project/internal/handlers/next_date"
	"go_final_project/internal/handlers/task"
	"go_final_project/internal/handlers/task_done"
	"go_final_project/internal/handlers/tasks"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/internal/db"
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

	authHandler := auth.NewHandler()
	nextdateHandler := next_date.NewHandler()
	taskHandler := task.NewHandler(dbInstance)
	tasksHandler := tasks.NewTasksHandler(dbInstance)
	taskDoneHandler := task_done.NewHandler(dbInstance)

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	http.HandleFunc("/api/signin", authHandler.Handle())
	http.HandleFunc("/api/nextdate", nextdateHandler.Handle())
	http.HandleFunc("/api/tasks", authHandler.Validate(tasksHandler.Handle()))
	http.HandleFunc("/api/task", authHandler.Validate(taskHandler.Handle()))
	http.HandleFunc("/api/task/done", authHandler.Validate(taskDoneHandler.Handle()))

	log.Printf("Сервер запущен на порту %d\n", port)
	log.Printf("Обслуживание файлов из каталога: %s\n", webDir)
	listenAddr := fmt.Sprintf("localhost:%d", port)
	log.Printf("Запуск сервера на %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
