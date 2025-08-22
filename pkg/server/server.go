package server

import (
	"go_final_project_v2/pkg/api"
	"go_final_project_v2/pkg/db"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DefaultPort = "7540"
	WebDir      = "./web"
)

// StartServer запускает веб-сервер
func StartServer(port string, webDir string) error {
	// Инициализируем базу данных
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	api.Init()
	// Создаем обработчик статического контента
	fs := CustomFileServer(webDir)
	http.Handle("/", fs)

	return http.ListenAndServe(":"+port, nil)
}

// CustomFileServer - расширенный обработчик файлов
func CustomFileServer(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Корректируем путь к файлу
		path := filepath.Join(root, r.URL.Path)

		// Если запрашивается корневой путь, возвращаем index.html
		if r.URL.Path == "/" {
			path = filepath.Join(root, "index.html")
		}

		// Проверяем существование файла
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// Служебная функция для обслуживания файла
		http.ServeFile(w, r, path)
	})
}
