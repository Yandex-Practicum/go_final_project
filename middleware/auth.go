package middleware

import (
	"net/http"

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
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			// ...

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
