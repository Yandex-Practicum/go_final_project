package main

import (
	"final_project/internal/authentification"
	addtask "final_project/internal/handlers/add_task"
	deletetask "final_project/internal/handlers/delete_task"
	donetask "final_project/internal/handlers/done_task"
	edittask "final_project/internal/handlers/edit_task"
	gettask "final_project/internal/handlers/get_task"
	gettasks "final_project/internal/handlers/get_tasks"
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

	addTaskHandler := addtask.New(repo)
	getTasksHandler := gettasks.New(repo)
	getTaskHandler := gettask.New(repo)
	editTaskHandler := edittask.New(repo)
	doneTaskHandler := donetask.New(repo)
	deleteTaskHandler := deletetask.New(repo)
	router.Handle("/*", http.FileServer(http.Dir("./web")))
	router.Post("/api/signin", authentification.Sign)
	router.Post("/api/task", authentification.Auth(addTaskHandler.AddTask))
	router.Get("/api/task", authentification.Auth(getTaskHandler.GetTaskHandler))
	router.Put("/api/task", authentification.Auth(editTaskHandler.EditTaskHandler))
	router.Delete("/api/task", authentification.Auth(deleteTaskHandler.DeleteTaskHandler))
	router.Post("/api/task/done", authentification.Auth(doneTaskHandler.DoneTaskHandler))
	router.Get("/api/nextdate", nextdate.NextDateHandler)
	router.Get("/api/tasks", authentification.Auth(getTasksHandler.GetTasksHandler))
	port := getPort()
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatal("ошибка создания сервера ", err)
	}
}
