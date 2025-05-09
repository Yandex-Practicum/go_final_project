package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func Run(r *chi.Mux) error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Error load env %w", err)
	}

	webDir := "./web"

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	if err := http.ListenAndServe(":"+port, r); err != nil {
		return err
	}

	return nil
}
