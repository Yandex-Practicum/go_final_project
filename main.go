package main

import (
	storage "go_final_project/database"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const defPort = "7070"

func main() {

	// Инициализируем ENV
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// Установка порта по умолчанию
	port := os.Getenv("TODO_PORT")

	if port != "" {
		port = defPort
	}

	// Запуск базы данных
	dbCon, err := storage.NewDb()
	if err != nil {
		log.Fatal(err)
	}

	defer dbCon.Db.Close()

	// Установка директории с фронтенд файлами
	fs := http.FileServer(http.Dir("./web"))

	// Обработка запросов к корневой директории
	http.Handle("/", fs)

	// Запуск сервера
	log.Printf("Server started at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
