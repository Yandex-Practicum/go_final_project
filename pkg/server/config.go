package server

import (
	"os"
)

const defaultPort = "7540"

// getPort возвращает порт из переменной окружения TODO_PORT
// или значение по умолчанию, если переменная не задана.
func getPort() string {
	port, found := os.LookupEnv("TODO_PORT")
	if !found || port == "" {
		return defaultPort
	}
	return port
}
