package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// generateJWT создает JWT-токен с полезной нагрузкой (хэш пароля)
func generateJWT(password string) (string, error) {
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default_secret_key" // Используем дефолтный ключ, если нет переменной окружения
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hashPassword(password),
		"exp":  time.Now().Add(8 * time.Hour).Unix(), // Токен действует 8 часов
	})

	return token.SignedString([]byte(secret))
}

// validateJWT проверяет валидность токена
func validateJWT(tokenString string) (bool, error) {
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default_secret_key"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, errors.New("invalid token")
	}

	// Проверяем соответствие хэша
	expectedHash := hashPassword(os.Getenv("TODO_PASSWORD"))
	if claims["hash"] != expectedHash {
		return false, errors.New("invalid token hash")
	}

	return true, nil
}

// hashPassword создает простейший хэш (в реальном проекте лучше использовать bcrypt)
func hashPassword(password string) string {
	return password + "_hashed"
}

// authMiddleware - Middleware для проверки аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
				return
			}

			token := cookie.Value
			storedHash, err := ValidateJWT(token)
			if err != nil || storedHash != pass {
				http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
