package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitAuth инициализирует секрет для JWT
func InitAuth() {
	password := os.Getenv("TODO_PASSWORD")
	if password != "" {
		jwtSecret = []byte(password)
	}
}

// GenerateToken создает JWT токен
func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": jwtSecret,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// ValidateToken проверяет JWT токен
func ValidateToken(tokenString string) bool {
	if jwtSecret == nil {
		return true
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

// AuthMiddleware проверяет аутентификацию
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if !ValidateToken(cookie.Value) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

// SigninHandler обрабатывает POST /api/signin
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		writeJSONError(w, "Аутентификация не настроена", http.StatusInternalServerError)
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if req.Password != password {
		writeJSONError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken()
	if err != nil {
		writeJSONError(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]string{"token": token}, http.StatusOK)
}
