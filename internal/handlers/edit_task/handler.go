package edittask

import "final_project/internal/repository"

type Handler struct {
	rep *repository.Repository
}

func New(rep *repository.Repository) *Handler {
	return &Handler{
		rep: rep,
	}
}
