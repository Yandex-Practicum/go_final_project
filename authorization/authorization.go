package authorization

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

// validateToken(tokenString) проводит валидации переданного jwt-токена
// возвращает true, если токен валиден, иначе false
func validateToken(tokenString string) (bool, error) {
	storedPassword := os.Getenv("TODO_PASSWORD")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(storedPassword), nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, fmt.Errorf("недействительный токен")
}

// Auth(next) реализует механизм middleware, проверяя аутентификации перед обработкой запроса
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtString string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtString = cookie.Value
			}
			var valid bool

			valid, err = validateToken(jwtString)
			if err != nil {
				valid = false
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
