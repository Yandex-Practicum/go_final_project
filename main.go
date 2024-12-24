package main

import (
	"os"
	"log"
	"strconv"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/FunnyFoXD/go_final_project/tests"
)

func main() {
	router := chi.NewRouter()
	port := getPort()

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	
	log.Printf("starting listen server on port %s", port)
	for err := http.ListenAndServe(":"+port, router); err != nil; {
		log.Fatalf("start server error: %s", err.Error())
	}
}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	return port
}
