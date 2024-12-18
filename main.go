package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Инициализация базы данных
	db, err := InitializeDatabase()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	// Регистрация обработчиков
	http.HandleFunc("/api/tasks", GetTasksHandler)
	http.HandleFunc("/api/task", EditTaskHandler)
	http.HandleFunc("/api/task/done", MarkTaskDoneHandler)

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
