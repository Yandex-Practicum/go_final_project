package main

import (
	"log"
	"os"

	"final_project/internal/db"
	"final_project/internal/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	err := db.Init(dbFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Db.Close()
	server.CreateServer()
}
