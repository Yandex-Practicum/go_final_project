package main

import (
	"fmt"
	"log"
	"os"

	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	// Получаем путь к БД из переменной окружения или используем по умолчанию
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	// Инициализация базы данных
	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Запуск веб-сервера
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	fmt.Println("Server is running...")
}
