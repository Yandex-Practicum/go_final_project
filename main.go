package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Sclea3/go_final_project/api"
	"github.com/Sclea3/go_final_project/db"
)

func main() {
	var err error
	api.DB, err = db.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}
	defer api.DB.Close()

	// Регистрируем маршруты для API.
	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.HandleFunc("/api/task/done", api.DoneHandler)

	// Файловый сервер для статики (директория "./web").
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// Определяем порт через переменную окружения или по умолчанию 7540.
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Println("Сервер запущен на порту", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
