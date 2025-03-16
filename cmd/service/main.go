package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/askoren1/go_final_project/internal/db"
	"github.com/askoren1/go_final_project/internal/handler"
	"github.com/askoren1/go_final_project/internal/repository"
)

func main() {

	//Инициализация
	dbConn := db.New() //Инициализируем базу данных
	defer db.Close(dbConn)
	repo := repository.New(dbConn) //Создаем репозиторий, который отвечает за взаимодействие с базой данных
	db.Migration(repo)             //Выполним миграцию базы данных

	handler := handler.New(repo) //Создадим обработчик HTTP-запросов, который будет использовать созданный репозиторий

	// Маршрутизация
	r := chi.NewRouter()                               //Создадим маршрутизатор
	r.Post("/api/task", handler.AddTask)               //Обработчик для добавления новой задачи
	r.Post("/api/task/done", handler.MarkTaskDone)     //Отмечает задачу как выполненную
	r.Get("/api/tasks", handler.GetList)               //Возвращает список задач
	r.Get("/api/nextdate", handler.NextDate)           // Возвращает следующую дату
	r.Get("/api/task", handler.GetTask)                //Получает информацию об одной задаче
	r.Put("/api/task", handler.UpdateTask)             //Обновляет информацию о задаче
	r.Delete("/api/task", handler.DeleteTask)          //Удаляет задачу
	r.Handle("/*", http.FileServer(http.Dir("./web"))) //Настраивает статический файловый сервер для реализации frontend

	//Запуск сервера
	port := os.Getenv("TODO_PORT") //Получаем порт из переменной окружения
	if port == "" {
		port = "8080" // Порт по умолчанию
	}
	if err := http.ListenAndServe(":"+port, r); err != nil { //Запускаем HTTP-сервер на указанном порту
		log.Fatal(err)
	}
}
