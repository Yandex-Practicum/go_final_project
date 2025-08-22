package main

import (
	"fmt"
	"go_final_project_v2/pkg/db"
	"log"
	"os"

	"go1f/pkg/server"
)

func main() {
	// Получаем порт из окружения или используем дефолтный
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = server.DefaultPort
	}

	// Получаем путь к БД из окружения или используем дефолтный
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем базу данных
	err := db.Init(dbFile)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Println("Ошибка закрытия базы данных:", err)
		}
	}()

	// Запускаем сервер
	fmt.Printf("Starting server on port %s\n", port)
	err = server.StartServer(port, server.WebDir)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
