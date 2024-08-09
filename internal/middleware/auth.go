package middleware

import (
	"net/http"
	"os"

	"go_final_project/internal/handlers"
)

type Auth struct {
	password string
}

func NewAuth() *Auth {
	pass := os.Getenv("TODO_PASSWORD")
	return &Auth{password: pass}
}

func (a *Auth) Handle(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(a.password) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			if !handlers.IsTokenValid(jwt) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
