package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go1f/pkg/api" // <-- укажи путь, соответствующий go.mod
)

func Run() error {
	port := 7540

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Регистрация API
	api.Init()

	// Файловый сервер (для фронта)
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	fmt.Printf("Server running on http://localhost:%d/\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
