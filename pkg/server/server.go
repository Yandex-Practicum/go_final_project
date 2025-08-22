package server

import (
	"fmt"
	"go_final_project/pkg/api"
	"log"
	"net/http"
	"os"
)

// Start принимает путь к директории с фронтендом и стартует HTTP-сервер.
// Если переменная окружения TODO_PORT задана, берёт порт из неё,
// иначе использует порт по умолчанию.
func Start(webDir string) {
	// читаем окружение
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// создаём файловый сервер для отдачи статичных файлов
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on http://localhost%s/, serving %s\n", addr, webDir)

	api.Init()

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
