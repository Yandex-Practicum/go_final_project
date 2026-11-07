package server

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

func Start() error {
	webDir := "./web"

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Сервер запущен на порту", port)

	return http.ListenAndServe(":"+port, nil)
}
