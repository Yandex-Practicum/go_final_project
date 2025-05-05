package main

import (
	"finalproject/pkg/db"
	"finalproject/pkg/server"
	"log"

	"github.com/go-chi/chi/v5"
)

func main() {
	_, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	if err := server.Run(r); err != nil {
		log.Printf("Could not start the server %v\n", err)
	}

}
