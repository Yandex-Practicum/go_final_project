package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

const webDir = "./web"

var DB *sql.DB
const formatDate = "20060102"

func main() {


	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	_, err = os.Stat(dbFile)
	if err != nil {
		createDb(DB)
	} else {
		log.Println("Database already exists")
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)
	http.Handle("/api/nextdate", http.HandlerFunc(ApiNextDate))
	http.HandleFunc("/api/task", Check)
	http.HandleFunc("/api/tasks", GetTasks)
	http.HandleFunc("/api/task/done", DoneTask)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}


