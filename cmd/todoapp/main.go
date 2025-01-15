package main

import (
	"github.com/ASHmanR17/go_final_project/internal/database"
	"github.com/ASHmanR17/go_final_project/internal/transport/httpserver"
)

func main() {
	// Проверим базу данных.
	database.Check()
	// Запускаем сервер
	httpserver.Start()
}
