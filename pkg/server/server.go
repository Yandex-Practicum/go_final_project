package server

import (
	"log"
	"net/http"
	"os"

	"github.com/seregey-schr/go_final_project/pkg/api"
	"github.com/seregey-schr/go_final_project/pkg/db"
)

func Start() {
	// Инициализация базы данных
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()
	log.Println("База данных успешно инициализирована")

	// Регистрация API
	api.Init()

	// Веб-сервер для фронтенда
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Сервер запущен на http://localhost:%s/ ...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
