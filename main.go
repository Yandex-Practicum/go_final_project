package main

import (
	"log"
	"os"

	"github.com/UTKIN-AI/go_final_project/pkg/api"
	"github.com/UTKIN-AI/go_final_project/pkg/db"
	"github.com/UTKIN-AI/go_final_project/pkg/server"
)

// main инициализирует базу данных, регистрирует API обработчики
// и запускает HTTP сервер.
func main() {
	dbFile := "scheduler.db"

	// Используем путь к базе данных из переменной окружения,
	// если она задана
	if envFile, found := os.LookupEnv("TODO_DBFILE"); found && envFile != "" {
		dbFile = envFile
	}

	err := db.Init(dbFile) 
	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}
	
	api.Init()

	server.Start()
	
}