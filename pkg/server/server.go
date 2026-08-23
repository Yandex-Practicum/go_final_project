package server

import (
	"log"
	"net/http"
)

func Start() {
	port := getPort()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("http server: %v", err)
	}
}
