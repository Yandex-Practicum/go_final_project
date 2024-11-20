package main

import (
	"final_project/internal/authentification"
	"final_project/internal/handlers"
	nextdate "final_project/internal/handlers/next_date"
	"final_project/internal/migration"

	"github.com/go-chi/chi"

	"log"
	"net/http"
)

var TimeFormat = "20060102"

func main() {

	repo := migration.Migration()

	router := chi.NewRouter()

	Handler := handlers.New(repo)

	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.Post("/api/signin", authentification.Sign)
	router.Post("/api/task", authentification.Auth(Handler.AddTask))
	router.Get("/api/task", authentification.Auth(Handler.GetTaskHandler))
	router.Put("/api/task", authentification.Auth(Handler.EditTaskHandler))
	router.Delete("/api/task", authentification.Auth(Handler.DeleteTaskHandler))
	router.Post("/api/task/done", authentification.Auth(Handler.DoneTaskHandler))
	router.Get("/api/nextdate", nextdate.NextDateHandler)
	router.Get("/api/tasks", authentification.Auth(Handler.GetTasksHandler))
	port := getPort()
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatal("ошибка создания сервера ", err)
	}
}
