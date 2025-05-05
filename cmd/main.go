package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load env ", err)
	}

	r := chi.NewRouter()

	webDir := "./web"

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	port := os.Getenv("TODO_PORT")
	http.ListenAndServe(":"+port, r)
}
