package main

import (
	"fmt"
	"go_final-project/tests"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webDir = "web"

func main() {
	// Порт для прослушивания
	port := tests.Port
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			log.Printf("Invalid TODO_PORT environment variable: %s. The default port is used: %d\n", envPort, port)
		} else {
			port = p
		}
	}

	// Обработчик файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Запуск сервера
	addrPort := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addrPort, nil)
	if err != nil {
		panic(err)
	}
}
