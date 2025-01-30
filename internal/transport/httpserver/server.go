package httpserver

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ASHmanR17/go_final_project/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type TaskServer struct {
	services service.TaskService
}

func NewTaskServer(s service.TaskService) *TaskServer {
	return &TaskServer{services: s}
}

// Serve запускает сервер
func (t *TaskServer) Serve() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %w", err)
	}
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

	//// Получаем пароль из переменной окружения
	//pass := os.Getenv("TODO_PASSWORD")
	//if len(pass) > 0 {
	////
	//}
	handler := newHTTPHandler(t.services) // handlers creating
	// Маршруты
	r.Get("/api/nextdate", handler.NextDate)
	r.Post("/api/task", auth(handler.AddTask))
	r.Get("/api/tasks", auth(handler.GetTasks))
	r.Get("/api/task", auth(handler.GetTask))
	r.Put("/api/task", auth(handler.EditTask))
	r.Post("/api/task/done", auth(handler.DoneTask))
	r.Delete("/api/task", auth(handler.DeleteTask))
	r.Post("/api/signin", handler.SignIn)

	// Запускаем сервер на порту 7540
	log.Printf("Serving files from %s on port %s", filesDir, port)

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":"+port, r))
}
