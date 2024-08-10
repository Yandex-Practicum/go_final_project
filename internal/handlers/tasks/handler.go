package tasks

import (
	"net/http"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

type TaskRepository interface {
	GetAllTasks() ([]*models.Task, error)
	GetAllTasksFilterByDate(date string) ([]*models.Task, error)
	GetAllTasksFilterByTitleOrComment(search string) ([]*models.Task, error)
}

type Handler struct {
	repository TaskRepository
}

func NewTasksHandler(repo TaskRepository) *Handler {
	return &Handler{repository: repo}
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
