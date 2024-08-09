package auth

import (
	"net/http"
	"os"
)

const (
	envPassword = "TODO_PASSWORD"
	secret      = "kjgsdfi632riusvd732fikd2--2!"
)

type Handler struct {
	password string
}

func NewHandler() *Handler {
	pass := os.Getenv("TODO_PASSWORD")
	return &Handler{password: pass}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.handleSign(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}
