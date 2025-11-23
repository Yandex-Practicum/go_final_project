package middleware

import (
	"crypto/sha256"
	"fmt"
	"go-task-manager-final_project/internal/api"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Auth - middleware-функция для проверки авторизации пользователя через JWT-токен.
// Параметр:
// next - обработчик HTTP-запроса, который будет вызван при успешной авторизации.
// Возвращает:
// http.HandlerFunc - обернутый обработчик с логикой авторизации.
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем значение пароля из переменной окружения TODO_PASSWORD.
		storedPassword := os.Getenv("TODO_PASSWORD")

		// Если пароль задан, выполняем проверку авторизации.
		if storedPassword != "" {
			// Пытаемся получить cookie с именем "token" из запроса.
			cookie, err := r.Cookie("token")
			if err != nil {
				// Если cookie отсутствует или возникла ошибка - возвращаем статус 401 (Неавторизован).
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
				return
			}

			// Получаем секрет для подписи JWT из переменной окружения JWT_SECRET.
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				api.WriteJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "JWT secret not configured",
				})
				return
			}
			secret := []byte(jwtSecret)

			// Парсим JWT-токен из значения cookie.
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				// Проверяем, что алгоритм подписи токена соответствует ожидаемому (HMAC).
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method %q", token.Header["alg"])
				}
				return secret, nil
			})

			// Если при парсинге токена произошла ошибка или токен недействителен - возвращаем ошибку.
			if err != nil || !token.Valid {
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "token expired or invalid",
				})
				return
			}

			// Извлекаем claims (данные) из токена.
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				// Если claims не соответствуют ожидаемому типу - возвращаем ошибку.
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid token: malformed claims",
				})
				return
			}

			// Вычисляем SHA-256 хэш текущего пароля из окружения.
			currentHash := sha256.Sum256([]byte(storedPassword))
			currentHashStr := fmt.Sprintf("%x", currentHash)

			// Сравниваем хэш пароля из токена с текущим хэшем пароля.
			// Если хэши не совпадают - токен недействителен.
			if claims["password_hash"] != currentHashStr {
				api.WriteJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid token: password changed",
				})
				return
			}

		}
		// Если все проверки прошли - передаём запрос дальше по цепочке обработчиков.
		next(w, r)
	})
}
