package handler

import "go_final_project/internal/repository"

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
