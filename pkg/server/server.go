package server

import (
	"fmt"
	"go1f/pkg/api"
	"log"
	"net/http"
	"os"
)

func StartServer() error {
	port := "7540"

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	api.Init()

	fs := http.FileServer(http.Dir("web"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}

		fs.ServeHTTP(w, r)
	})

	http.Handle("/css/", fs)
	http.Handle("/js/", fs)
	http.Handle("/favicon.ico", fs)

	log.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}
