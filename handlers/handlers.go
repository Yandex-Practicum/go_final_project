package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dekatei/go_final_project/task"
	_ "modernc.org/sqlite"
)

// 3ий шаг
func NextDateHandler(w http.ResponseWriter, req *http.Request) {
	nowStr := req.URL.Query().Get("now")
	if nowStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("now missing"))
		err := errors.New("пропущено время")
		fmt.Println(err)
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong time value"))
		fmt.Println(err)
	}

	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	nextDate, err := task.NextDate(now, date, repeat)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}

// 4ый шаг
func TaskHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch req.Method {
	case http.MethodGet:
		task.GetTask(w, req, db)
	case http.MethodPost:
		task.AddTask(w, req, db)
	case http.MethodPut:
		task.UpdateTask(w, req, db)
	case http.MethodDelete:
		task.DeleteTask(w, req, db)
	default:
		http.Error(w, fmt.Sprintf("Сервер не поддерживает %s запросы", req.Method),
			http.StatusMethodNotAllowed)
		return
	}
}

// 7ой шаг
func TaskDoneDeleteHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch req.Method {
	case http.MethodPost:
		task.DoneTask(w, req, db)
	case http.MethodDelete:
		task.DeleteTask(w, req, db)
	default:
		http.Error(w, fmt.Sprintf("Сервер не поддерживает %s запросы", req.Method),
			http.StatusMethodNotAllowed)
		return
	}
}

func GetTasksHandler(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch req.Method {
	case http.MethodGet:
		task.GetTasks(w, req, db)
	default:
		http.Error(w, fmt.Sprintf("Сервер не поддерживает %s запросы", req.Method),
			http.StatusMethodNotAllowed)
		return
	}
}
