package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"main.go/database"
	"main.go/parsedate"
	"main.go/tasks"
	_ "modernc.org/sqlite"
)

const (
	defPort   = "7540"
	webDir    = "./web"
	defDBFile = "./scheduler.db"
)

func main() {
	var db *sql.DB
	var err error
	db, err = database.InitDatabase()
	if err != nil {
		log.Fatalf("Bad DB: %v\n", err)
	}
	defer db.Close()

	envik := godotenv.Load()
	if envik != nil {
		log.Fatal("Error env file")
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defPort
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", parsedate.NextDateHandler)
	//http.HandleFunc("/api/task", tasks.AddTaskHandler(db))
	http.HandleFunc("/api/tasks", tasks.GetTasksHandler(db))
	http.HandleFunc("/api/task/done", tasks.DoneMarkHandler(db))
	http.HandleFunc("/api/signin", parsedate.SignHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			tasks.AddTaskHandler(db)(w, r)
		case http.MethodGet:
			tasks.GetTaskHandler(db)(w, r)
		case http.MethodPut:
			tasks.UpdateTaskHandler(db)(w, r)
		case http.MethodDelete:
			tasks.DeleteTaskHandler(db)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(tasks.ErrorResponse{Error: "Method Not Allowed"})
		}
	})

	log.Printf("Server on: %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server start error: %v", err)
	}
}
