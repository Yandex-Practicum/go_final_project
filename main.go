package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VesuvyX/go_final_project/database"
	"github.com/VesuvyX/go_final_project/handlers"
	"github.com/VesuvyX/go_final_project/models"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

const defaultPort = "7540"

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}

func main() {
	port := getPort()
	const webDir = "web"
	//const webDir = "./web"

	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("ошибка подключения к бд: %v", err)
	}
	defer db.Close()

	models.SetDB(db)

	handler := chi.NewRouter()
	fs := http.FileServer(http.Dir(webDir))

	handler.Mount("/", fs)
	handler.Get("/api/nextdate", handlers.NextDateGETHandler)
	handler.Post("/api/task", models.TaskAddPOST)
	handler.Get("/api/tasks", models.TasksShowGET)
	handler.Get("/api/task", models.ReadTaskByIdGET)
	handler.Put("/api/task", models.TaskUpdatePUT)
	handler.Post("/api/task/done", models.TaskDonePOST) // выполнение задачи
	handler.Delete("/api/task/done", models.TaskDELETE) // удаление задачи

	fmt.Printf("Запуск сервера на порту %s ...\n\n", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		panic(err)
	}
}
