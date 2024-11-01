package main

import (
	"fmt"
	"net/http"
	"todo-server/internal/authtentification"
	"todo-server/internal/database"
	"todo-server/internal/handlers"
	"todo-server/internal/settings"
	"todo-server/internal/store"

	"github.com/go-chi/chi"
)

func main() {
	db := database.InitializeDatabase()
	defer db.Close()

	store := store.NewStorage(db)

	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.Post("/api/task", authtentification.Authentification(handlers.HandlePostGetPutRequests(store)))
	r.Get("/api/tasks", authtentification.Authentification(handlers.HandleTasksGet(store)))
	r.Get("/api/task", handlers.HandlePostGetPutRequests(store))
	r.Put("/api/task", handlers.HandlePostGetPutRequests(store))
	r.Post("/api/task/done", authtentification.Authentification(handlers.HandleTaskDone(store)))
	r.Delete("/api/task", handlers.HandlePostGetPutRequests(store))
	r.Post("/api/signin", authtentification.HandleSiginingIn)

	environment := settings.GetEnv()
	err := http.ListenAndServe(":"+environment.Port, r)
	if err != nil {
		fmt.Println("Не удалось запустить сервер:\n", err)
	}
}
