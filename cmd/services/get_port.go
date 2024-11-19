package main

import (
	"log"
	"os"
	"strconv"
)

func getPort() string {
	var port string
	s := os.Getenv("TODO_PORT")
	if s != "" {
		_, err := strconv.Atoi(s)
		if err != nil {
			port = "7540"
			log.Println("неверно указан номер порта в переменной окружения: сервер будет использовать порт по умолчанию(:7540)")
			return port
		}
		return s
	} else {
		port = "7540"
		log.Println("не указан номер порта в пременной окружения: сервер будет использовать порт по умолчанию(:7540)")
		return port
	}
}
