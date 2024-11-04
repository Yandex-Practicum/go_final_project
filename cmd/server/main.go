package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/internal/api"
	"go_final_project/internal/db"

	"github.com/joho/godotenv" // Импортируем библиотеку для загрузки .env
)

func main() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить файл .env. Используются переменные окружения системы.")
	}

	// Получаем порт из переменной окружения или устанавливаем значение по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
		log.Println("Запускается сервер на порту по умолчанию:", port)
	} else {
		log.Println("Запускается сервер на порту из переменной окружения:", port)
	}

	// Инициализируем базу данных
	database, err := db.SetupDB()
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.Close() // Закрываем базу данных при завершении работы сервера

	// Устанавливаем директорию для файлов фронтенда
	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Добавляем обработчики для API
	http.HandleFunc("/api/nextdate", api.NextDateHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос к /api/task: метод=%s", r.Method) // Логируем метод запроса
		switch r.Method {
		case http.MethodGet:
			api.GetTaskHandler(w, r) // Обработчик GET-запроса для получения задачи по ID
		case http.MethodPut:
			api.UpdateTaskHandler(w, r) // Обработчик PUT-запроса для обновления задачи
		case http.MethodPost:
			api.AddTaskHandler(w, r) // Обработчик POST-запроса для добавления новой задачи
		case http.MethodDelete:
			api.DeleteTaskHandler(w, r) // Обработчик DELETE-запроса для удаления задачи
		default:
			http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос к /api/tasks: метод=%s", r.Method) // Логируем метод запроса
		api.GetUpcomingTasks(database)(w, r)
	}) // Обработчик для списка задач

	http.HandleFunc("/api/task/done", api.CompleteTaskHandler) // Обработчик для выполнения задачи

	// Логируем и запускаем сервер
	log.Printf("Запуск сервера на порту %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
