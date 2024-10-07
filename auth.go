package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims структура для JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var appPassword string // Глобальная переменная для хранения пароля

// init инициализирует глобальные переменные
func init() {
	// Получаем значение пароля из переменной окружения один раз при старте приложения
	appPassword = os.Getenv("TODO_PASSWORD")
}

// authMiddleware — middleware для проверки JWT токена
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appPassword != "" { // Если пароль установлен, проверяем аутентификацию
			cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Authentication required", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			tokenStr := cookie.Value
			claims := &Claims{}

			// Здесь предполагается, что jwtKey уже определён где-то в другом месте кода
			tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			if !tkn.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Если пароль не установлен, продолжаем без проверки
		next(w, r)
	}
}

// generateToken создает новый JWT токен
func generateToken() (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		Username: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Здесь предполагается, что jwtKey уже определён где-то в другом месте кода
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
