package gettasks

import "final_project/internal/repository"

type Handler struct {
	rep *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		rep: repo,
	}
}
