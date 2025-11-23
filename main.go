package main

import (
	"go-task-manager-final_project/internal/db"
	"go-task-manager-final_project/internal/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// В main инициализируем соединение с базой данных, обеспечиваем его корректное закрытие и запускаем HTTP-сервер для обработки запросов.
func main() {
	// Загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found, using system environment variables")
		} else {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	// Открываем соединения с БД и, при необходимости, создаем схему
	db, err := db.Init("")
	if err != nil {
		log.Printf("failed to initialize database: %v", err)
	}
	// Обеспечиваем закрытие соединения с БД при завершении работы программы (даже в случае паники или ошибки).
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database connection: %v", closeErr)
		}
	}()

	// Запускаем сервер
	err = server.StartServer(db)
	if err != nil {
		log.Printf("failed to start server: %v", err)
		return
	}
}
