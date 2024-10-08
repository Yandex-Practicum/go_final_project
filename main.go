package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	db, err := setupDB()
	if err != nil {
		log.Fatalf("Ошибка настройки БД: %v", err)
	}
	defer db.Close()

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// API рабы с аутентификацией
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", authMidW(makeHandler(taskHandler, db)))
	http.HandleFunc("/api/tasks", authMidW(makeHandler(tasksHandler, db)))
	http.HandleFunc("/api/task/done", authMidW(makeHandler(taskDoneHandler, db)))
	http.HandleFunc("/api/signin", makeHandler(signInHandler, db))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Обработчик для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r, db)
	case http.MethodGet:
		tasksHandler(w, r, db)
	case http.MethodPut:
		editTaskHandler(w, r, db)
	case http.MethodDelete:
		deleteTaskHandler(w, r, db)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
