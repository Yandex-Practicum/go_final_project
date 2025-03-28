package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dekatei/go_final_project/base"
	"github.com/dekatei/go_final_project/handlers"
)

const webDir = "./web"

func main() {
	envDBFILE := os.Getenv("TODO_DBFILE")
	db, err := base.CreateDB(envDBFILE)
	// Подключаемся к БД
	if err != nil {
		log.Printf("1")
		fmt.Println(err)
		return

	}
	defer db.Close()
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc(`/api/nextdate`, handlers.NextDateHandler)
	http.HandleFunc(`/api/task`, func(w http.ResponseWriter, req *http.Request) {
		handlers.TaskHandler(w, req, db)
	})
	http.HandleFunc(`/api/tasks`, func(w http.ResponseWriter, req *http.Request) {
		handlers.GetTasksHandler(w, req, db)
	})
	http.HandleFunc(`/api/task/done`, func(w http.ResponseWriter, req *http.Request) {
		handlers.TaskDoneDeleteHandler(w, req, db)
	})
	envPort := os.Getenv("TODO_PORT")
	if envPort == "" {
		envPort = "7540"
	}

	fmt.Printf("Запускаем сервер на порту: %s", envPort)

	err = http.ListenAndServe(":"+envPort, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Завершаем работу")
}
