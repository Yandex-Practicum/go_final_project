package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/qwsnxnjene/go_final_project/authorization"
	"github.com/qwsnxnjene/go_final_project/db"
	"github.com/qwsnxnjene/go_final_project/handlers"
	"github.com/qwsnxnjene/go_final_project/tests"
)

// startServer() обеспечивает начало работы сервера, создание БД и настройку API
func startServer() {
	ports := ":" + strconv.Itoa(tests.Port)
	webDir := "./web"

	if p := os.Getenv("PORT"); p != "" {
		ports = ":" + p
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("ошибка загрузки .env файла")
	}

	db.CreateDB()
	defer db.Database.Close()

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", authorization.Auth(handlers.TaskHandler))
	http.HandleFunc("/api/tasks", authorization.Auth(handlers.TasksHandler))
	http.HandleFunc("/api/task/done", authorization.Auth(handlers.TaskDoneHandler))
	http.HandleFunc("/api/signin", handlers.SignInHandler)
	fmt.Printf("Сервер запущен на порте %v\n", ports)
	err = http.ListenAndServe(ports, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("can't start a server: %v", err))
		return
	}

}

func main() {
	startServer()
}
