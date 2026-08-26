package server

import (
	"fmt"
	"net/http"
)

func Start() error {
	port := getPort()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}
