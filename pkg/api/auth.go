package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {

			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			writeJson(w, map[string]string{"error": "authentication required"}, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(pass), nil
		})
		if err != nil {
			writeJson(w, map[string]string{"error": "invalid token"}, http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			writeJson(w, map[string]string{"error": "invalid token"}, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeJson(w, map[string]string{"error": "invalid token claims"}, http.StatusUnauthorized)
			return
		}
		tokenHash, ok := claims["password_hash"].(string)
		if !ok {
			writeJson(w, map[string]string{"error": "invalid token claims"}, http.StatusUnauthorized)
			return
		}

		hash := sha256.Sum256([]byte(pass))
		currentHash := hex.EncodeToString(hash[:])

		if tokenHash != currentHash {
			writeJson(w, map[string]string{"error": "token is no longer valid"}, http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
