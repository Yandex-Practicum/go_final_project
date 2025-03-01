package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"todo_restapi/internal/myfunctions"
	"todo_restapi/internal/storage"
)

// go test -run ^TestApp$ ./tests
// go test -run ^TestDB$ ./tests

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("no .env file found: %v\n", err)
	}
}

func main() {

	t, e := myfunctions.NextDate(time.Now(), "", "")
	fmt.Println(t, e)

	port, exists := os.LookupEnv("TODO_PORT")
	if !exists {
		port = "8080"
	}

	storagePath, exists := os.LookupEnv("TODO_DBFILE")
	if !exists {
		storagePath = "scheduler.db"
	}

	database, err := storage.OpenStorage(storagePath)
	if err != nil {
		log.Fatalf("OpenStorage: %v\n", err)
	}

	_ = database

	router := chi.NewRouter()

	router.Get("/", func(write http.ResponseWriter, request *http.Request) {
		http.ServeFile(write, request, "web/index.html")
	})

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("web"))))

	fmt.Printf("Server is running on port%s...\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("server run error: %v\n", err)
	}
}
