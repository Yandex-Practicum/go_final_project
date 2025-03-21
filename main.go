package main

import (
	"log"
	"net/http"

	"go_final_project/api"
	"go_final_project/db"
)

func main() {
	const port = ":7540"
	const webDir = "./web"

	// Инициализируем базу данных
	database, err := db.InitDB(false)
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}
	defer database.Close()

	log.Println("База данных успешно инициализирована.")

	// Настраиваем файловый сервер для обслуживания статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Регистрация обработчиков для API
	http.HandleFunc("/api/nextdate", api.NextDateHandler) // Используем обработчик из nextdate.go
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		api.GetTasksHandler(w, r, database) // Обработчик для получения задач
	})
	http.HandleFunc("/api/addtask", func(w http.ResponseWriter, r *http.Request) {
		api.AddTaskHandler(w, r, database) // Обработчик для добавления задач
	})

	log.Printf("Сервер запущен на http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
