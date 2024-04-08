package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var cookieToken string
			cookie, err := r.Cookie("token")
			if err == nil {
				cookieToken = cookie.Value
			}
			jwtInstance := jwt.New(jwt.SigningMethodHS256)
			token, err := jwtInstance.SignedString([]byte(pass))

			if cookieToken != token {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
