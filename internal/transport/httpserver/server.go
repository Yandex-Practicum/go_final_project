package httpserver

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Start запускает сервер
func Start() {

	// Настройка роутера
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Указываем директорию, из которой нужно обслуживать файлы
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web"))

	// Настраиваем сервер для обслуживания файлов из указанной директории
	FileServer(r, "/", filesDir)

	// Если существует переменная окружения TODO_PORT, сервер при старте должен слушать порт со значением этой переменной.
	defaultPort := "7540"
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// подключаемся к БД
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Маршруты
	r.Get("/api/nextdate", handleNextDate)
	r.Post("/api/task", handleAddTask(db))
	r.Get("/api/tasks", handleGetTasks(db))
	r.Get("/api/task", handleGetTask(db))
	r.Put("/api/task", handleEditTask(db))
	r.Post("/api/task/done", handleDoneTask(db))
	r.Delete("/api/task", handleDeleteTask(db))

	// Запускаем сервер на порту 7540
	log.Printf("Serving files from %s on port %s", filesDir, port)

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":"+port, r))
}
