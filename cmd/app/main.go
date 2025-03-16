package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/internal/constants"

	"go_final_project/internal/db"
	"go_final_project/internal/handlers"
	"go_final_project/internal/repository"
	"go_final_project/internal/service"
)

func main() {
	port := constants.DefaultPort
	if portStr := os.Getenv(constants.EnvPort); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	dbInstance, err := db.InitDatabase()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer dbInstance.Close()

	taskRepository := repository.NewTaskRepository(dbInstance)
	taskService := service.NewTaskService(taskRepository)
	authService := service.NewAuthService()

	authHandler := handlers.NewAuthHandler(authService)
	nextDateHandler := handlers.NewNextDateHandler()
	taskHandler := handlers.NewTaskHandler(taskService)
	tasksHandler := handlers.NewTasksHandler(taskService)
	taskDoneHandler := handlers.NewTaskDoneHandler(taskService)

	fs := http.FileServer(http.Dir(constants.WebDir))
	http.Handle("/", fs)

	http.HandleFunc("/api/signin", authHandler.Handle())
	http.HandleFunc("/api/nextdate", nextDateHandler.Handle())
	http.HandleFunc("/api/tasks", authHandler.Validate(tasksHandler.Handle()))
	http.HandleFunc("/api/task", authHandler.Validate(taskHandler.Handle()))
	http.HandleFunc("/api/task/done", authHandler.Validate(taskDoneHandler.Handle()))

	log.Printf("Сервер запущен на порту %d\n", port)
	log.Printf("Обслуживание файлов из каталога: %s\n", constants.WebDir)
	listenAddr := fmt.Sprintf("localhost:%d", port)
	log.Printf("Запуск сервера на %s\n", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
