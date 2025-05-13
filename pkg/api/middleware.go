package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := godotenv.Load(); err != nil {
			return
		}

		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtToken string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				errHandler(w, "", err, http.StatusBadRequest)
				return
			}

			if !token.Valid {
				errHandler(w, "Authentification required", errors.New(""), http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if claims["user_password"] != os.Getenv("TODO_PASSWORD") {
					errHandler(w, "Invalid token payload", nil, http.StatusUnauthorized)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
