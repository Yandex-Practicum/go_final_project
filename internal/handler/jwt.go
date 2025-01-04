package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret ключ для подписания JWT
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Signin обрабатывает вход пользователя
func Signin(w http.ResponseWriter, r *http.Request) {
	// Чтение JSON из тела запроса
	var payload struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Проверка пароля
	expectedPassword := os.Getenv("TODO_PASSWORD")
	if payload.Password != expectedPassword {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	// // Создание JWT-токена
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"hash": fmt.Sprintf("%x", sha256.Sum256([]byte(expectedPassword))),
	// 	"exp":  time.Now().Add(8 * time.Hour).Unix(),
	// })
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Ответ с токеном
	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
