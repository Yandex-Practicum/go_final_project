package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/LTVgreater5CPi/go_final_project/tasks_service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	appPassword = os.Getenv("TODO_PASSWORD")
	if appPassword == "" {
		log.Println("Переменная TODO_PASSWORD не установлена. Аутентификация отключена.")
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
	http.HandleFunc("/api/nextdate", tasks_service.NextDateHandler)
	http.HandleFunc("/api/task", authMidW(tasks_service.MakeHandler(taskHandler, db)))
	http.HandleFunc("/api/tasks", authMidW(tasks_service.MakeHandler(tasks_service.TasksHandler, db)))
	http.HandleFunc("/api/task/done", authMidW(tasks_service.MakeHandler(tasks_service.TaskDoneHandler, db)))
	http.HandleFunc("/api/signin", tasks_service.MakeHandler(signInHandler, db))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Обработчик для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		tasks_service.AddTaskHandler(w, r, db)
	case http.MethodGet:
		tasks_service.TasksHandler(w, r, db)
	case http.MethodPut:
		tasks_service.EditTaskHandler(w, r, db)
	case http.MethodDelete:
		tasks_service.DeleteTaskHandler(w, r, db)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
