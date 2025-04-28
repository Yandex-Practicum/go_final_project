package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

// HTTP-обработчик для аутентификации пользователей
func signinHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, "Недопустимый формат запроса")
		return
	}

	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" {
		writeError(w, "Аутентификация не настроена")
		return
	}

	if creds.Password != envPassword {
		writeError(w, "Неверный пароль")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pass_hash": hashPassword(envPassword),
		"exp":       time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		writeError(w, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(8 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	writeJSON(w, map[string]string{"token": tokenString})
}

func hashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}
