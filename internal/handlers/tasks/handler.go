package tasks

import (
	"database/sql"
	"net/http"

	"go_final_project/internal/utils"
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
			http.Error(w, utils.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}
