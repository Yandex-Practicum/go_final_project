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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var tokenString string

			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}
			valid := validateToken(tokenString, pass)

			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func validateToken(tokenString string, password string) bool {
	secret := getJWTSecret()

	if secret == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return false
	}

	if !token.Valid {
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

func getJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}