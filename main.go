package main

import (
	"go1f/pkg/db"
	"go1f/pkg/server"
	"log"
)

func main() {

	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := server.StartServer(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
