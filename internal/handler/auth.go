package handler

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"go_final_project/internal/model"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req model.SignInRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" {
		http.Error(w, "Аутентификация не настроена", http.StatusInternalServerError)
		return
	}

	if req.Password != envPassword {
		response := model.SignInResponse{Error: "Неверный пароль"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": envPassword,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(envPassword))
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	response := model.SignInResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envPassword := os.Getenv("TODO_PASSWORD")
		if len(envPassword) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(envPassword), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["hash"] != envPassword {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
