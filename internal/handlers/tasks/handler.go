package tasks

import (
	"database/sql"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func NewTasksHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetTasks(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}
