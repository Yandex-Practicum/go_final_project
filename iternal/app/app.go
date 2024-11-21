package app

import (
	"Go/iternal/database"

	"github.com/go-chi/chi"

	"fmt"
	"net/http"
)

const (
	webDir = "./web/"
)

func Run() {
	_, err := database.CreateDB()
	if err != nil {
		panic(err)
	}
	router := chi.NewRouter()
	fmt.Println("Запускаем сервер!")

	router.Handle("/*", http.StripPrefix("/web", http.FileServer(http.Dir(webDir))))

	err = http.ListenAndServe(":7540", router)
	if err != nil {
		panic(err)
	}
}
