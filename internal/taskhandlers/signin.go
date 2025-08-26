package taskhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go_final_project/internal/middleware"
)

// SignInRequest структура для пароля
type SignInRequest struct {
	Password string `json:"password"`
}

// SignInHandler обрабатывает вход пользователя
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Проверяем пароль
	storedPassword := os.Getenv("TODO_PASSWORD")
	if storedPassword == "" {
		http.Error(w, `{"error":"Server password is not set"}`, http.StatusInternalServerError)
		return
	}
	if req.Password != storedPassword {
		http.Error(w, `{"error":"Invalid password"}`, http.StatusUnauthorized)
		return
	}

	// Генерируем JWT-токен
	token, err := middleware.GenerateJWT(req.Password)
	if err != nil {
		http.Error(w, `{"error":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Отправляем токен в ответе
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Println("User authenticated successfully")
}
