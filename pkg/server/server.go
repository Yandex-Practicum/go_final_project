package server

import (
	"log"
	"net/http"
)

const WebDir = "./web"

func Run(port string) error {
	http.Handle("/", http.FileServer(http.Dir(WebDir)))

	log.Printf("Сервер запущен на http://localhost:%s/\n", port)

	err := http.ListenAndServe(":"+port, nil)

	return err
}
