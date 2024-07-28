package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/AlexJudin/go_final_project/config"
)

type AuthMW struct {
	cfg *config.Сonfig
}

func New(c *config.Сonfig) AuthMW {
	return AuthMW{cfg: c}
}

func (a *AuthMW) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := a.cfg.Password
		if len(pass) > 0 {
			var (
				signedToken string
				password    string
			)
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				signedToken = cookie.Value
			}
			jwtToken, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
				// секретный ключ для всех токенов одинаковый, поэтому просто возвращаем его
				return password, nil
			})
			if err != nil {
				fmt.Errorf("Failed to parse token: %s\n", err)
				return
			}

			if !jwtToken.Valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
