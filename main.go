package main

import (
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

const (
	port    = ":7540"
	webDir  = "./web"
	dateFmt = "20060102"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "Not all parameters passed", http.StatusBadRequest)
		return
	}

	now, err := time.Parse(dateFmt, nowStr)
	if err != nil {
		http.Error(w, "Incorrect format now", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}

func main() {
	db, err := InitDB()
	if err != nil {
		log.Fatal("Error creating database:", err)
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleCreateTask(w, r, db)
		} else if r.Method == http.MethodGet {
			if r.URL.Query().Has("id") {
				handleGetTaskByID(w, r, db)
			} else {
				http.Error(w, `{"error":"Не указан ID задачи"}`, http.StatusBadRequest)
			}
		} else if r.Method == http.MethodPut {
			handleUpdateTask(w, r, db)
		} else if r.Method == http.MethodDelete {
			handleDeleteTask(w, r, db)
		} else {
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGetTasks(w, r, db)
		} else {
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleTaskDone(w, r, db)
		} else {
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Server running http://localhost%s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
