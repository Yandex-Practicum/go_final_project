package main

import (
	"go_final_project/internal/server"
	"go_final_project/pkg/db"
	"log"
)

func main() {
	logs := log.Default()

	conn, err := db.Init("./scheduler.db")
	if err != nil {
		logs.Fatal(err)
	}

	defer conn.Close()

	server.Init("7540", "./web", conn, logs)
}
