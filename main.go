package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AndreiCoder404/go_final_project/pkg/db"
	"github.com/AndreiCoder404/go_final_project/pkg/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if os.Getenv("GO_TEST") == "1" {
		log.SetOutput(io.Discard)
	} else {
		log.SetFlags(0)
	}

	fmt.Println("Init: Opening database file: scheduler.db")
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}
	fmt.Println("Database initialized successfully")

	err = server.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
