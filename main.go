package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	dbA "go_final_project/myLib/dataBase"
	httpH "go_final_project/myLib/httpH"

	"github.com/go-chi/chi/v5"
)

const webDir = "./web/" // директория с файлами

func main() {

	// Подключение БД - Создание БД
	// Отключение от БД
	//
	sqlDB, dbCl, err := dbA.CheckCreateDB()
	if err != nil {
		if dbCl != nil {
			dbCl()
		}
		log.Fatal(err)
		os.Exit(1)
	}
	err = sqlDB.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	// Обработчики URL
	//
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	r.Get("/api/nextdate", httpH.NextDateH)    // Формирование новой даты
	r.Post("/api/task", httpH.AddTasksH)       // Добавление задач
	r.Get("/api/tasks", httpH.ReadTasksH)      // Чтение задач
	r.Get("/api/task", httpH.GetTasksH)        // Чтение параметров задачи по его ID
	r.Put("/api/task", httpH.SaveTasksH)       // Сохранение данных задачи по его ID
	r.Post("/api/task/done", httpH.DoneTasksH) // Завершение задачи
	r.Delete("/api/task", httpH.DelTasksH)     // Удаление задачи

	// Запуск сервера
	//
	fmt.Println("Запуск сервера.")

	err = http.ListenAndServe(":"+os.Getenv("HTTP_PORT"), r)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		log.Fatal(err)
		os.Exit(1)
	}

}
