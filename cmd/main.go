package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/AlexJudin/go_final_project/api"
	"github.com/AlexJudin/go_final_project/database"
	"github.com/AlexJudin/go_final_project/usecases"
)

// @title Пользовательская документация API
// @description Учебный сервис (Яндекс Практикум)
// @termsOfService spdante@mail.ru
// @contact.name Alexey Yudin
// @contact.email spdante@mail.ru
// @version 1.0.0
// @host http://localhost:7540
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Error connect to database: %+v", err)
	}
	defer db.Close()

	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		port = "7540"
	}

	// init usecases
	taskUC := usecases.NewTaskUsecase(db)
	taskHandler := api.NewTaskHandler(taskUC)

	webDir := "./web"
	r := chi.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", taskHandler.GetNextDate)

	serverAddress := fmt.Sprintf("localhost:%s", port)
	log.Println("Listening on " + serverAddress)
	if err = http.ListenAndServe(serverAddress, r); err != nil {
		log.Panicf("Start server error: %+v", err.Error())
	}
}
