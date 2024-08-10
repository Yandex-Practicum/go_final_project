package handlers

import (
	"encoding/json"
	"go_final_project/internal/models"
	"go_final_project/internal/utils"
	"net/http"
)

type AuthService interface {
	IsAuthEnabled() bool
	IsTokenValid(token string) bool
	SignIn(password string) (string, error)
}

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{svc: service}
}

func (h *AuthHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.handleSign(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func (h *AuthHandler) Validate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.svc.IsAuthEnabled() {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			if !h.svc.IsTokenValid(jwt) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func (h *AuthHandler) handleSign(w http.ResponseWriter, r *http.Request) {
	var signDTO models.SignDTO
	err := json.NewDecoder(r.Body).Decode(&signDTO)
	if err != nil {
		utils.RespondWithError(w, utils.ErrInvalidJson)
		return
	}

	token, err := h.svc.SignIn(signDTO.Password)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	response := models.SignResponseDTO{Token: token}
	utils.Respond(w, response)
}
