package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Инициализация БД
	if err := initDB(); err != nil {
		log.Fatal("DB init error:", err)
	}
	defer DB.Close()

	http.HandleFunc("/api/task/done", TaskDoneHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", GetTasksHandler)
	http.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		nowStr := r.FormValue("now")
		dateStr := r.FormValue("date")
		repeat := r.FormValue("repeat")

		now, err := time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "Неверный формат даты 'now'", http.StatusBadRequest)
			return
		}

		result, err := NextDate(now, dateStr, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(result))
	})
	// Запуск сервера
	port := getPort()
	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Printf("Server started at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return "7540"
}
