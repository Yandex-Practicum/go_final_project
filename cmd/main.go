package main

import (
	"go_final_project/cmd/route"
	"go_final_project/config"
	"log"
	"net/http"
	"os"
)

func main() {
	// Получение значения порта из переменной окружения TODO_PORT
	port := os.Getenv("TODO_PORT")

	// Инициализация базы данных
	db := config.InitializeDatabase()
	defer db.Close()

	// Настройка сервера
	router := route.SetupRouter()
	log.Printf("Сервер запущен на порту: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
