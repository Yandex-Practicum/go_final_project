package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if appPassword != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Неавторизован", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Bad req", http.StatusBadRequest)
				return
			}

			tokenStr := cookie.Value
			claims := &Claims{}

			tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "Неавторизован", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Bad req", http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				http.Error(w, "Неавторизован", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func generateToken() (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &Claims{
		Username: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
