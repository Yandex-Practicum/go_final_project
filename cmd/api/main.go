package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"todo-list/pkg/api"
	"todo-list/pkg/db"
)

func main() {

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if !filepath.IsAbs(dbFile) {

		execPath, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		if filepath.Base(execPath) == "api" {
			execPath = filepath.Dir(execPath)
			execPath = filepath.Dir(execPath)
		}

		dbFile = filepath.Join(execPath, dbFile)
	}

	log.Printf("Файл базы данных: %s", dbFile)

	if err := db.Init(dbFile); err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}
	defer db.Close()
	log.Println("База данных успешно подключена")

	api.Init()
	log.Println("API обработчики")

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	webDir := filepath.Join(path, "web")

	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = filepath.Join(path, "..", "web")
	}

	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = "./web"
		log.Printf("Путь: %s", webDir)
	}

	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Printf("Папка web не найдена по пути: %s", webDir)
	} else {
		log.Printf("Обслуживаю статические файлы из: %s", webDir)
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %s", port)
	log.Printf("Откройте http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
