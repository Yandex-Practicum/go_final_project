package main

import (
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
	"os"
)

const (
	defaultPort = "7540"
	WebDir      = "./web"
)

func main() {
	// определение пути к файлу (со звездочкой)
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.InitDB(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы %v", err)
	}
	// со звездочкой
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	//инициализация обработчиков
	api.InitAPI()

	err = server.Run(port)

	if err != nil {
		log.Fatal(err)
	}
}
