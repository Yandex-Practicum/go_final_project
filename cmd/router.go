package main

import (
	"github.com/gorilla/mux"
	"github.com/kmakasheva/todo-list-project/handlers"
	"github.com/kmakasheva/todo-list-project/middleware"
	"net/http"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/task", middleware.Auth(handlers.PostTaskHandler)).Methods("POST")
	r.HandleFunc("/api/task", middleware.Auth(handlers.GetTaskHandler)).Methods("GET")
	r.HandleFunc("/api/task", middleware.Auth(handlers.UpdateTaskHandler)).Methods("PUT")
	r.HandleFunc("/api/task", middleware.Auth(handlers.DeleteTaskHandler)).Methods("DELETE")
	r.HandleFunc("/api/tasks", middleware.Auth(handlers.GetTasksHandler))
	r.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	r.HandleFunc("/api/task/done", middleware.Auth(handlers.DoneTaskHandler)).Methods("POST")
	r.HandleFunc("/api/signin", handlers.SignInHandler).Methods("POST")

	webDir := "./web/"
	fileServer := http.FileServer(http.Dir(webDir))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	return r
}
