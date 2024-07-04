package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"main.go/database"
	"main.go/handlers"
	"main.go/tests"

	_ "modernc.org/sqlite"
)

func main() {
	// Проверка базы данных
	// Создание базы данных и таблицы, если не существует
	db, err := database.СheckDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	webDir := "./web"

	r := chi.NewRouter()
	fs := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fs)

	// Регистрация обработчиков
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.Post("/api/task", handlers.HandleTask(db))
	r.Get("/api/task", handlers.HandleTaskID(db))
	r.Put("/api/task", handlers.HandleTask(db))
	r.Delete("/api/task", handlers.HandleTask(db))
	r.Post("/api/task/done", handlers.HandleTaskDone(db))
	r.Get("/api/tasks", handlers.HandleTask(db))

	port := fmt.Sprintf(":%d", tests.Port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}
}
