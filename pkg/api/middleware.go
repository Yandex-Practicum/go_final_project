package api

import (
	"github.com/JKasus/go_final_project/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, err := config.NewConfig()
		if err != nil {
			http.Error(w, "error loading config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if cfg.Password != "" {
			var tokenString string
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}

			var valid bool
			if tokenString != "" {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.Password), nil
				})
				if err == nil && token.Valid {
					valid = true
				}
			}

			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
