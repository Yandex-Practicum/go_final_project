package handlers

import (
	"database/sql"
	"net/http"
)

func RegisterHandlers(db *sql.DB) {
	http.HandleFunc("/api/nextdate", NextDateHandler(db))
	http.HandleFunc("/api/task", TaskHandler(db))
	http.HandleFunc("/api/tasks", TasksHandler(db))
	http.HandleFunc("/api/task/done", DoneTaskHandler(db))
}
