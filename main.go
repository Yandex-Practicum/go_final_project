package main

import (
	"log"
	"net/http"
	"os"
)

// Структура для хранения задачи
type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func main() {
	// Инициализация базы данных
	db = initDB()
	defer db.Close()

	// Настройка маршрутов
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/task", addTaskHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)

	// Запуск сервера
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	log.Printf("Сервер запущен на http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
