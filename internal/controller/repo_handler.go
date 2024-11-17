package controller

import "pwd/internal/repository"

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
