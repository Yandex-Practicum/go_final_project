package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AlexJudin/go_final_project/api"
	"github.com/AlexJudin/go_final_project/config"
	"github.com/AlexJudin/go_final_project/repository"
	"github.com/AlexJudin/go_final_project/usecases"
)

// @title Пользовательская документация API
// @description Итоговый проект (Яндекс Практикум)
// @termsOfService spdante@mail.ru
// @contact.name Alexey Yudin
// @contact.email spdante@mail.ru
// @version 1.0.0
// @host http://localhost:7540
// @BasePath /api/
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewDB(cfg.DBFile)
	if err != nil {
		log.Fatalf("Error connect to repository: %+v", err)
	}
	defer db.Close()

	// init repository
	repo := repository.NewNewRepository(db)

	// init usecases
	taskUC := usecases.NewTaskUsecase(repo)
	taskHandler := api.NewTaskHandler(taskUC)

	webDir := "./web"
	r := chi.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", taskHandler.GetNextDate)
	r.Post("/api/task", taskHandler.CreateTask)
	r.Get("/api/tasks", taskHandler.GetTasks)
	r.Get("/api/task", taskHandler.GetTask)
	r.Put("/api/task", taskHandler.UpdateTask)
	r.Post("/api/task/done", taskHandler.MakeTaskDone)
	r.Delete("/api/task", taskHandler.DeleteTask)

	serverAddress := fmt.Sprintf("localhost:%s", cfg.Port)
	log.Println("Listening on " + serverAddress)
	if err = http.ListenAndServe(serverAddress, r); err != nil {
		log.Panicf("Start server error: %+v", err.Error())
	}
}
