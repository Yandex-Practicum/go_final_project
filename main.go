package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	port   int
	webDir string = "./web"
)

func main() {
	portFromEnv := os.Getenv("TODO_PORT")
	if portFromEnv != "" {
		p, err := strconv.Atoi(portFromEnv)
		if err != nil {
			log.Fatalf("Failed to parse TODO_PORT: %v", err)
		}
		port = p
	} else {
		port = 7540
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", NextDateHandler)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on port %d\n", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
