package main

import (
	"net/http"
	"pwd/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	r.Get("/", handlers.MainHandle)

	http.ListenAndServe(":7540", nil)
}
