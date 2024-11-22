package app

import (
	"Go/iternal/database"

	"fmt"
	"net/http"
)

const (
	webDir = "./web"
)

func Run() {
	_, err := database.CreateDB()
	if err != nil {
		panic(err)
	}
	fmt.Println("Запускаем сервер!")

	http.Handle("/*", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
