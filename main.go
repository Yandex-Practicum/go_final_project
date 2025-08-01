package main

import (
	"log"
	"os"

	"github.com/ZhakiyaK/final_project/pkg/db"
	"github.com/ZhakiyaK/final_project/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {
	logger := log.New(os.Stdout, "APP: ", log.LstdFlags)

	// Подключение к БД с использованием логгера
	if _, err := os.Stat("scheduler.db"); os.IsNotExist(err) {
		logger.Fatal("DB file not found")
	}

	if err := db.Init("scheduler.db"); err != nil {
		logger.Fatal("DB init failed: ", err)
	}
	defer db.Close()

	// Запуск сервера с передачей логгера
	srv := server.NewServer(logger)

	if err := srv.Start(); err != nil {
		logger.Fatal("Error starting server: ", err)
	}
}
