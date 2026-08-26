package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	pass := os.Getenv("TODO_PASSWORD") // проверяем наличие пароля, читаем один раз при регистрации middleware

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pass == "" {
			next(w, r) // авторизация не требуется
			return
		}

		var tokenString string

		cookie, err := r.Cookie("token")
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "authentication required",
			})
			return
		}
		tokenString = cookie.Value

		valid := validateToken(tokenString, pass)
		if !valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "authentication required",
			})
			return
		}
		next(w, r)
	})
}

func validateToken(tokenString string, password string) bool {
	if password == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(password), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	hashFromToken, ok := claims["hash"].(string)
	if !ok {
		return false
	}

	hash := sha256.Sum256([]byte(password))
	currentHash := hex.EncodeToString(hash[:])

	return hashFromToken == currentHash
}
