package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go_final-project/internal/db"
	"go_final-project/internal/handlers"
	"log"
	"net/http"
	"os"
	"strconv"
)

const Port = 7540

func main() {
	// Порт для прослушивания
	port := Port
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			log.Printf("Invalid TODO_PORT environment variable: %s. The default port is used: %d\n", envPort, port)
		} else {
			port = p
		}
	}

	// Инициализация базы данных
	database, closeDB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error to initialize db: %v", err)
	}
	defer closeDB()

	// Обработчик файлов
	http.Handle("/", http.FileServer(http.Dir("web")))

	http.HandleFunc("/tasks", handlers.GetTasksHandler(database))
	http.HandleFunc("/addTask", handlers.AddTaskHandler(database))

	// Запуск сервера
	addrPort := fmt.Sprintf(":%d", port)
	err = http.ListenAndServe(addrPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
