package main

import (
	"go_final_project/pkg/db"
	. "go_final_project/pkg/server"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	// Инициализация БД
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
		// ИЛИ panic("Ошибка инициализации БД: " + err.Error())
	}
	defer db.Close()

	// Запуск сервера
	StartServer()

}
