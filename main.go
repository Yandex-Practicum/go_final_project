package main

import (
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close() // Закрытие базы данных после завершения работы программы

	log.Println("Запуск сервера...")
	server.Run() // Теперь запускаем сервер после инициализации базы данных
}
