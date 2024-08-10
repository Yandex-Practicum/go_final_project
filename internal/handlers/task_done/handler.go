package task_done

import (
	"net/http"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

type TaskRepository interface {
	GetTaskByID(ID int64) (*models.Task, error)
	UpdateTaskDate(taskID int64, date string) error
	DeleteTaskByID(ID int64) error
}

type Handler struct {
	repository TaskRepository
}

func NewHandler(repository TaskRepository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.handlePostTaskDone(w, r)
		default:
			http.Error(w, utils.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}
