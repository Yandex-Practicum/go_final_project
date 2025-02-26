package main

import (
	"fmt"
	"net/http"
	"os"

	"go_final_project/handlers"
	"go_final_project/sqlite"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {

	port := os.Getenv("TODO_PORT")
	todoDB := os.Getenv("TODO_DBFILE")

	if port == "" {
		port = ":7540"
	}

	// Проверить существует ли файл БД. Если его нет, то создать БД.
	dbURL, err := sqlite.FindOrCreateDB(todoDB)
	if err != nil {
		fmt.Println("Ошибка с базой данных ", err)
	}

	// Подключаемся к БД
	db, err := sqlite.InitDB(dbURL)
	defer db.Close()
	if err != nil {
		fmt.Println("Ошибка инициализации БД ", err)
	}

	sqlite.TodoStorage = sqlite.NewStorage(db)

	// Создаём роутер
	r := chi.NewRouter()

	// Запускаем Web интерфейс
	r.Handle("/*", http.FileServer(http.Dir("./web")))

	// Выводим значение новой даты
	r.Get("/api/nextdate", handlers.GetNextDateHandler)

	// Работаем с одной задачей
	r.Post("/api/task", handlers.PostOneTaskHandler)
	r.Get("/api/task", handlers.GetOneTaskHandler)
	r.Put("/api/task", handlers.PutOneTaskHandler)
	r.Post("/api/task/done", handlers.DoneOneTaskHandler)
	r.Delete("/api/task", handlers.DeleteOneTaskHandler)

	// Работа с группой задач
	r.Get("/api/tasks", handlers.GetTasksHandler)

	// Запускаем сервер
	fmt.Printf("Сервер TODO запущен! Порт %s.\n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s.\n", err.Error())
	}
}
