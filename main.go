package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	// Инициализация базы данных
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	// Получение порта из переменной окружения или значения по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Настройка маршрутов
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
	http.HandleFunc("/api/task/done", markTaskDoneHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTaskByIDHandler(w, r)
		} else if r.Method == http.MethodPut {
			updateTaskHandler(w, r)
		} else if r.Method == http.MethodPost {
			addTaskHandler(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTaskHandler(w, r) // обработка DELETE-запроса для удаления задачи
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	// Запуск сервера
	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
