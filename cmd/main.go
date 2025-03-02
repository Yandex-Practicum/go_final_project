package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	const webDir = "./web"

	fs := http.FileServer(http.Dir(webDir))

	r.Handle("/*", fs)

	log.Println("считай уже стартанул")
	err := http.ListenAndServe(":7540", r)
	if err != nil {
		panic(err)
	}
}
