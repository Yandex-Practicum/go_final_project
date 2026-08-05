package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"proj/pkg/api/service"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := service.AutherENV()
		if len(pass) <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		var claims service.Claims

		cookie, err := r.Cookie("token")
		if err != nil {
			service.SendErrorJSON(w, "Authentification required", http.StatusUnauthorized)
			return
		}
		jwtToken, err := jwt.ParseWithClaims(cookie.Value, &claims, func(t *jwt.Token) (any, error) {
			return []byte(service.Secret), nil
		})
		if err != nil || !jwtToken.Valid {
			service.SendErrorJSON(w, "Authentification required", http.StatusUnauthorized)

			return
		}
		hashedPassword := sha256.Sum256([]byte(pass))
		hashStringPassword := hex.EncodeToString(hashedPassword[:])

		if claims.PasswordHash != hashStringPassword {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

	})
}
