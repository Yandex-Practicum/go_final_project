package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/database"
	"go_final_project/handlers"
)

func main() {
	port := "7540"
	if val, ok := os.LookupEnv("TODO_PORT"); ok {
		port = val
	}

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Инициализация базы данных
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Регистрация обработчиков
	handlers.RegisterHandlers(db)

	log.Printf("Starting server on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
