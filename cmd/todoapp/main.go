package main

import (
	"log"

	"github.com/ASHmanR17/go_final_project/internal/database"
	"github.com/ASHmanR17/go_final_project/internal/service"
	"github.com/ASHmanR17/go_final_project/internal/transport/httpserver"
)

func main() {

	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := database.NewTaskStore(db)
	services := service.NewTaskService(*store)
	taskServer := httpserver.NewTaskServer(*services)

	// Запускаем сервер
	taskServer.Serve()

}
