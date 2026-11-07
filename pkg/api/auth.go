package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 8 * time.Hour

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func createToken(password string) (string, error) {
	claims := Claims{
		PasswordHash: hashPassword(password),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(password))
}

func validateToken(tokenString string, password string) bool {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("неверный алгоритм подписи")
			}

			return []byte(password), nil
		},
	)

	if err != nil || !token.Valid {
		return false
	}

	expectedHash := hashPassword(password)

	return subtle.ConstantTimeCompare(
		[]byte(claims.PasswordHash),
		[]byte(expectedHash),
	) == 1
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if configuredPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value, configuredPassword) {
			http.Error(
				w,
				"Требуется авторизация",
				http.StatusUnauthorized,
			)
			return
		}

		next(w, r)
	}
}
