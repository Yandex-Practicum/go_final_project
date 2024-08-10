package auth

import (
	"go_final_project/internal/utils"
	"net/http"
	"os"
)

type Handler struct {
	password string
}

func NewHandler() *Handler {
	pass := os.Getenv(utils.EnvPassword)
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
