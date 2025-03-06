package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sclea3/go_final_project/db"
	"github.com/Sclea3/go_final_project/handlers"
)

func main() {
	// Определяем путь к базе данных относительно исполняемого файла.
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbPath := filepath.Join(filepath.Dir(exePath), "scheduler.db")

	// Для разработки удаляем существующую базу (если есть)
	if _, err := os.Stat(dbPath); err == nil {
		log.Printf("Удаляем существующую базу данных: %s", dbPath)
		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Ошибка удаления базы: %v", err)
		}
	}

	// Инициализируем базу (создаём таблицу и индекс)
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Регистрируем API-обработчики
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.TaskHandler(w, r, database)
	})
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		handlers.ListTasksHandler(w, r, database)
	})
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		handlers.TaskDoneHandler(w, r, database)
	})

	// Статические файлы фронтенда
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("./web")))
	http.Handle("/index.html", http.FileServer(http.Dir("./web")))
	http.Handle("/login.html", http.FileServer(http.Dir("./web")))

	port := "7540"
	log.Printf("Сервер стартанул на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
