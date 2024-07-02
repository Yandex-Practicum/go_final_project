package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Memonagi/go_final_project/database"
	"github.com/Memonagi/go_final_project/handler"
	"github.com/Memonagi/go_final_project/tests"
	"github.com/go-chi/chi/v5"
)

const (
	webDir = "./web"
)

func main() {
	// получение значения переменной окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	// подключение к базе данных
	db, err := database.CheckDatabase()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// создание маршрутизатора и обработка запросов
	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// вычисление следующей даты
	r.Get("/api/nextdate", handler.GetNextDate)
	// добавление задачи в БД
	r.MethodFunc(http.MethodPost, "/api/task", handler.PostAddTask)
	// получение списка задач
	r.MethodFunc(http.MethodGet, "/api/tasks", handler.GetAddTasks)
	// получение задачи по ее идентификатору
	r.MethodFunc(http.MethodGet, "/api/task", handler.GetTaskId)
	// редактирование задачи
	r.MethodFunc(http.MethodPut, "/api/task", handler.UpdateTaskId)
	// выполнение задачи
	r.MethodFunc(http.MethodPost, "/api/task/done", handler.TaskDone)
	// удаление задачи
	r.MethodFunc(http.MethodDelete, "/api/task", handler.DeleteTask)

	// запуск сервера
	log.Printf("запуск веб-сервера на порту %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		fmt.Println(err)
	}
}
