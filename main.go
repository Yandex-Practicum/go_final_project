package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go_final_project-main/internal/config"
	"go_final_project-main/internal/database"
	"go_final_project-main/internal/handlers"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %+v", err)
	}

	db, err := sqlx.Open("sqlite3", cfg.DbFile)
	if err != nil {
		log.Fatalf("error opening db: %+v", err)
	}
	defer db.Close()

	err = database.CheckDb(db, cfg)
	if err != nil {
		log.Fatalf("error checking db: %+v", err)
	}

	webDir := "./web/"
	port := "7540"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/tasks", handlers.GetTasksHandler(db))
	http.HandleFunc("/api/task/done", handlers.MarkTaskDoneHandler(db))
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddTaskHandler(db)(w, r)
		case http.MethodGet:
			handlers.GetTaskHandler(db)(w, r)
		case http.MethodPut:
			handlers.EditTaskHandler(db)(w, r)
		case http.MethodDelete:
			handlers.DeleteTaskHandler(db)(w, r)

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Запуск веб-сервера на порту %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
