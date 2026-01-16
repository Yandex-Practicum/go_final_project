package server

import (
	"go_final_project/pkg/api"
	"log"
	"net/http"
)

// StartServer запускает сервер с обработчиками запросов
func StartServer() {
	webDir := "./web" // Директория с файлами фронта

	// Api эндпоинты
	api.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if urlPath == "/" {
			urlPath = "/index.html"
		}
		pathToFile := webDir + urlPath
		// Логгирую запрос
		log.Println("По запросу", urlPath, "возвращен файл", pathToFile)
		http.ServeFile(w, r, pathToFile)
	})

	// Запуск сервера
	port := "7540"

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	} else {
		log.Println("Сервер запущен на порту", port)
	}
}
