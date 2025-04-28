package server

import (
	"log"
	"net/http"
	"os"
)

// InitServer запускает сервер
func InitServer() {
	http.Handle("/", http.FileServer(http.Dir("./web"))) // обработчик
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	addr := ":" + port
	log.Printf("Server listening on %s", addr)
	err := http.ListenAndServe(addr, nil) // запуск
	if err != nil {
		panic(err)
	}
}
