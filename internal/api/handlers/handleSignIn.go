package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go-task-manager-final_project/internal/api"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// signInRequest - структура для приёма данных из запроса на авторизацию.
// Содержит единственное поле:
// Password - пароль пользователя в виде строки (сериализуется как "password" в JSON).
type signInRequest struct {
	Password string `json:"password"`
}

// signInHandler - обработчик HTTP-запроса на авторизацию пользователя.
// Ожидает JSON с полем "password", проверяет пароль и возвращает JWT-токен при успехе.
// Параметры:
// w - объект http.ResponseWriter для отправки ответа клиенту.
// r - объект *http.Request с данными запроса.
func handleSignIn(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON из тела запроса в структуру signInRequest.
	// Если декодирование не удалось, возвращаем ошибку 400 (Bad Request).
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON format",
		})
		return
	}

	// Проверяем, что поле Password не пустое.
	// Если пароль пустой, возвращаем ошибку 400 (Bad Request).
	if req.Password == "" {
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "password cannot be empty",
		})
		return
	}

	// Получаем мастер-пароль из переменной окружения TODO_PASSWORD.
	// Если переменная не задана, возвращаем ошибку 500 (Internal Server Error).
	masterPass := os.Getenv("TODO_PASSWORD")
	if masterPass == "" {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "TODO_PASSWORD environment variable is not set",
		})
		return
	}

	// Сравниваем пароль из запроса с мастер-паролем.
	// Если пароли не совпадают, возвращаем ошибку 401 (Unauthorized).
	if req.Password != masterPass {
		api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "incorrect password",
		})
		return
	}

	// Получаем секрет для подписи JWT из переменной окружения JWT_SECRET.
	// Если переменная не задана, используем значение по умолчанию "bobr_kurwa".
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "JWT secret not configured",
		})
		return
	}
	secret := []byte(jwtSecret)

	// Вычисляем хэш пароля с помощью алгоритма SHA-256.
	hash := sha256.Sum256([]byte(req.Password))

	// Формируем claims (полезную нагрузку) JWT-токена:
	// - "authenticated": флаг успешной аутентификации (true).
	// - "exp": время истечения токена (текущее время + 8 часов).
	// - "iss": идентификатор сервера-издателя токена.
	// - "password_hash": шестнадцатеричное представление хэша пароля.
	claims := jwt.MapClaims{
		"authenticated": true,
		"exp":           time.Now().Add(time.Hour * 8).Unix(),
		"iss":           "go-task-manager-final_project",
		"password_hash": fmt.Sprintf("%x", hash),
	}

	// Создаём JWT-токен с указанными claims и алгоритмом подписи HS256.
	// Подписываем токен секретом и получаем его строковое представление.
	// При ошибке подписи возвращаем ошибку 500 (Internal Server Error).
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to generate JWT token",
		})
		return
	}

	// Возвращаем успешный ответ 200 (OK) с JWT-токеном в поле "token".
	api.WriteJSON(w, http.StatusOK, map[string]string{
		"token": signedToken,
	})

}
