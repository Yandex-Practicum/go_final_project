package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"final-project-bronner/go/db"
	"final-project-bronner/go/handlers"
)

func main() {
	connection := db.CreateConnectionDB()
	defer connection.Close()

	db := db.GetDB(connection)
	handler := handlers.NewHandler(db)
	router := http.NewServeMux()

	router.Handle("/", http.FileServer(http.Dir("../web")))
	router.HandleFunc("/api/nextdate", handler.NextDateHandler)
	router.HandleFunc("POST /api/task", handler.CreateTask)
	router.HandleFunc("POST /api/task/done", handler.MakeTaskHasDone)
	router.HandleFunc("GET /api/tasks", handler.GetAllTasks)
	router.HandleFunc("GET /api/task", handler.GetTask)
	router.HandleFunc("PUT /api/task", handler.UpdateTask)
	router.HandleFunc("DELETE /api/task", handler.DeleteTask)

	serverPort := os.Getenv("TODO_PORT")

	log.Printf("Сервер запущен на http://localhost:%s", serverPort)

	err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), router)
	if err != nil {
		log.Fatal(err)
	}
}
