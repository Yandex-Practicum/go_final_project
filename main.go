package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/Arukasnobnes/go_final_project/database"
	"github.com/Arukasnobnes/go_final_project/handlers"
	"github.com/Arukasnobnes/go_final_project/server"
	"github.com/Arukasnobnes/go_final_project/storage"
)

func main() {
	dbFilePath, err := db.GetDBFilePath()
	if err != nil {
		log.Fatalf("Can't get DB file path: %v", err)
	}

	db.InitDB(dbFilePath)

	s := storage.NewStorage(db.DB)
	h := handlers.NewHandler(s)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		fmt.Println("TODO_PORT not finding, use port 7540")
		port = "7540"
	}

	server.InitHandlers(h)

	log.Printf("Start server: port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
