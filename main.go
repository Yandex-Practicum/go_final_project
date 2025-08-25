package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

var port = os.Getenv("TODO_PORT")
var dbPath = os.Getenv("TODO_DBFILE")

func main() {

	//Переменная окружения для БД
	if dbPath == "" {
		dbPath = "pkg/db/scheduler.db"
	}

	//Соединение с БД
	if err := db.Init(dbPath); err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}

	//Закрываем соединение с БД
	defer db.Close()

	//Регистрация Эндпоинтов
	api.Init()

	//Настройка Статики
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	//Запуск сервера
	if port == "" {
		port = "7540"
	}

	addr := ":" + port

	log.Printf("Сервер запущен на %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
