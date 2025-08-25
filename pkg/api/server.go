package api

import (
	"fmt"
	"net/http"
)

func Run() error {

	//Настроим порт
	port := 7540

	//Статика
	http.Handle("/", http.FileServer(http.Dir("web")))

	//Запуск сервера
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
