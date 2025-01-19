package handlers

import "final-project-bronner/go/db"

type Handler struct {
	db *db.DB
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{
		db: db,
	}
}
