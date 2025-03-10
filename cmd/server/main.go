package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go_final-project/internal/db"
	"go_final-project/internal/handlers"
	"log"
	"net/http"
	"os"
	"strconv"
)

const Port = 7540

func main() {
	// Порт для прослушивания
	port := Port
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			log.Printf("Invalid TODO_PORT environment variable: %s. The default port is used: %d\n", envPort, port)
		} else {
			port = p
		}
	}

	// Инициализация базы данных
	database, closeDB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error to initialize db: %v", err)
	}
	defer closeDB()

	// Обработчик файлов
	http.Handle("/", http.FileServer(http.Dir("web")))

	// API маршруты
	http.HandleFunc("/api/tasks", handlers.GetTasksHandler(database))
	http.HandleFunc("/api/task", handlers.TaskHandler(database))
	http.HandleFunc("/api/task/done", handlers.MarkTaskDoneHandler(database))
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler(database))

	// Запуск сервера
	addrPort := fmt.Sprintf(":%d", port)
	log.Printf("Server started on http://localhost: %s\n", addrPort)
	err = http.ListenAndServe(addrPort, nil)
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", addrPort, err)
	}
}
