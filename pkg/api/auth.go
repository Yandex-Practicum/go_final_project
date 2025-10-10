package api

import (
	"encoding/json"
	"fmt"
	"go_final_project/pkg/auth"
	"net/http"
	"os"
	"time"
)

// SingInRequest структура для запроса входа
type SingInRequest struct {
	Password string `json:"password"`
}

// SingInResponse структура для ответа входа
type SingInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// SingInHandler обработчик для аутентификации
func SingInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, SingInResponse{Error: "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	var req SingInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, SingInResponse{Error: fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	// Проверяем пароль
	envPassword := getEnvPassword()
	if envPassword == "" {
		writeJSON(w, SingInResponse{Error: "Аутентификация не настроена"}, http.StatusInternalServerError)
		return
	}

	if req.Password != envPassword {
		writeJSON(w, SingInResponse{Error: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := auth.GenerateToken()
	if err != nil {
		writeJSON(w, SingInResponse{Error: fmt.Sprintf("Ошибка генерации токена: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(8 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, SingInResponse{Token: token}, http.StatusOK)
}

// AuthMiddleware middleware для проверки аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Если аутентификация отключена, пропускаем
		if !auth.IsAuthEnabled() {
			next(w, r)
			return
		}

		// Получаем токен из куки
		var token string
		cookie, err := r.Cookie("token")
		if err == nil {
			token = cookie.Value
		}

		// Если куки нет, пробуем получить токен из заголовка Authorization
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			}
		}

		valid, err := auth.ValidateToken(token)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}

// Вспомогательная функция для получения пароля из env
func getEnvPassword() string {
	return os.Getenv("TODO_PASSWORD")
}
