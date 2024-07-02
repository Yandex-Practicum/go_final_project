package server

import (
	"fmt"
	"main/tests"
	"os"
	"strconv"
)

func GetPort() int {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		return tests.Port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintf("Неверный порт: %s", portStr))
	}
	return port
}
