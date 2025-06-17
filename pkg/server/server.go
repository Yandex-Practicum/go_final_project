package server

import (
	"go_final_project/pkg/api"
	"log"
	"net/http"
	"os"
)

func Run() {

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Сервер запущен и слушает на порту %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
