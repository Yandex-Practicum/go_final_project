package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error":"Ошибка в запросе"}`, http.StatusBadRequest)
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")

	if request.Password != storedPassword {
		http.Error(w, `{"error":"Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["password_hash"] = storedPassword
	claims["exp"] = time.Now().Add(8 * time.Hour).Unix()

	secret := []byte("your-secret-key")
	signedToken, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, `{"error":"Ошибка создания токена"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, `{"error":"Необходима авторизация"}`, http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // Секретный ключ для верификации
		})

		if err != nil || claims["password_hash"] != pass {
			http.Error(w, `{"error":"Неверный токен или пароль"}`, http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
