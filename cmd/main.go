package main

import (
	"go_final_project/cmd/route"
	"go_final_project/config"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	port := os.Getenv("TODO_PORT")

	// Инициализация базы данных
	if err := config.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	// Настройка сервера
	router := route.SetupRouter()
	log.Printf("Сервер запущен на порту: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
