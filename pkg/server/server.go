package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func Run(r *chi.Mux) error {
	webDir := "./web"

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	port := "7540"
	if p := os.Getenv("TODO_PORT"); p != "" {
		port = p
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		return err
	}

	return nil
}
