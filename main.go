package main

import (
	"log"
	"net/http"
	"os"

	//"path/filepath"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"

	database "github.com/sandrinasava/go_final_project/db"
	"github.com/sandrinasava/go_final_project/handlers"
)

func main() {

	// Загружаю переменные окружения
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	// Подключаюсь к бд
	db := database.ConnectDB()
	defer db.Close()

	http.HandleFunc("/api/tasks", handlers.TasksHandler(db))

	http.HandleFunc("/api/task/done", handlers.TaskDoneHandler(db))

	http.HandleFunc("/api/task", handlers.TaskHandler(db))

	http.HandleFunc("/api/nextdate", handlers.NextDateHandle) // добавление обработчика в глобальный маршрутизатор

	fileServer := http.FileServer(http.Dir("./web")) // обработчик для директории WEB
	http.Handle("/", fileServer)                     // добавление обработчика в глобальный маршрутизатор

	// Получение значения TODO_PORT
	port := os.Getenv("TODO_PORT")
	if port == "" {
		log.Fatal("Переменная окружения TODO_PORT не установлена")
	}

	log.Println("Запуск сервера на порту:", port)
	// Запуск HTTP-сервера
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка запуска сервера: %v", err)
	}

}
