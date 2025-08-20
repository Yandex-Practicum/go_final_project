package main

import (
	"Final/pkg/db"
	"Final/pkg/server"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	fmt.Println("TODO_PASSWORD =", os.Getenv("TODO_PASSWORD"))

	dbFile := "scheduler.db"
	if dbEnvFile := os.Getenv("TODO_DBFILE"); dbEnvFile != "" {
		dbFile = dbEnvFile
	}

	if err := db.Init(dbFile); err != nil {
		fmt.Printf("Ошибка инициализации базы данных: %v\n", err)
		os.Exit(1)
	}
	defer db.DB.Close()

	server.Run()
}

//
