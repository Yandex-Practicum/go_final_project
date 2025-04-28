package main

import (
	"final_sprint/pkg/api"
	"final_sprint/pkg/db"
	"final_sprint/pkg/server"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Определение и инициализация базы данных через переменную окружения
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = filepath.Join(".", "scheduler.db")
	}
	if err := db.Init(dbPath); err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	api.Init()
	// Запуск сервера
	server.InitServer()

}
