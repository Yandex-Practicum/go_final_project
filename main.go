package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/qavaleria/go_final_project/database"
	"github.com/qavaleria/go_final_project/handlers"
	"github.com/qavaleria/go_final_project/tests"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Получаем порт из переменной окружения или используем значение по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	// Добавляем обработчик API для вычисления следующей даты
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.MethodFunc(http.MethodGet, "/api/task", handlers.HandleGetTask(db))
	r.MethodFunc(http.MethodPut, "/api/task", handlers.HandleUpdateTask(db))
	r.MethodFunc(http.MethodDelete, "/api/task", handlers.HandleDeleteTask(db))
	r.MethodFunc(http.MethodPost, "/api/task", handlers.HandleAddTask(db))
	r.MethodFunc(http.MethodGet, "/api/tasks", handlers.HandleGetTasks(db))
	r.MethodFunc(http.MethodPost, "/api/task/done", handlers.HandleMarkTaskDone(db))

	// Запускаем сервер
	log.Printf("Server is listening on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
