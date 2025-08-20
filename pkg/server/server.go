package server

import (
	"fmt"
	"net/http"
	"os"

	"Final/pkg/api"
)

func Run() {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}
	api.Init()

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
