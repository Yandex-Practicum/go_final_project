package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret ключ для проверки JWT
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware проверяет JWT токен
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			// Получаем токен из куки
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Проверяем JWT
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				// Проверяем, что используется правильный метод подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Проверка истечения
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}
			}

			// Проверяем хэш пароля
			hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
			if claims["hash"] != hashedPassword {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
		}

		// Если проверка прошла, передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
