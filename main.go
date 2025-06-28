package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	// определяем путь к файлу БД: из окружения TODO_DBFILE или по умолчанию
	dbFile := os.Getenv("TODO_DBFILE") // вернёт "" если не задано :contentReference[oaicite:10]{index=10}
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// инициализируем БД
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Database init error: %v", err)
	}

	// запускаем HTTP-сервер
	server.Start("web")
}
