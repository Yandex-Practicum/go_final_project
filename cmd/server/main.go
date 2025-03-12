package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go_final-project/internal/auth"
	"go_final-project/internal/db"
	"go_final-project/internal/handlers"
	"go_final-project/internal/logic"
	"log"
	"net/http"
	"os"
	"strconv"
)

const Port = 7540

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		log.Fatal("TODO_PASSWORD is not set in .env")
	}
	auth.GetSecretKey(password)

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
	http.HandleFunc("/api/signin", handlers.SignInHandler)
	http.HandleFunc("/api/tasks", auth.AuthMiddleware(handlers.GetTasksHandler(database)))
	http.HandleFunc("/api/task", auth.AuthMiddleware(handlers.TaskHandler(database)))
	http.HandleFunc("/api/task/done", auth.AuthMiddleware(handlers.MarkTaskDoneHandler(database)))
	http.HandleFunc("/api/nextdate", logic.NextDateHandler(database))

	// Запуск сервера
	addrPort := fmt.Sprintf(":%d", port)
	log.Printf("Server started on http://localhost: %s\n", addrPort)
	err = http.ListenAndServe(addrPort, nil)
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", addrPort, err)
	}
}
