package main

import (
	"github.com/Yandex-Practicum/go_final_project/pkg/db"
	"github.com/Yandex-Practicum/go_final_project/pkg/server"
	"log"
)

func main() {
	database, err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	server.Start(database)
}
