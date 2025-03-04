package middleware

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getSecretKey загружает ключ из переменной окружения перед каждым вызовом
func getSecretKey() []byte {
	return []byte(os.Getenv("TODO_PASSWORD"))
}

// GenerateJWT создает токен для аутентификации
func GenerateJWT(password string) (string, error) {
	secretKey := getSecretKey()
	if len(secretKey) == 0 {
		return "", errors.New("TODO_PASSWORD is not set")
	}

	// Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": password,
		"exp":           time.Now().Add(8 * time.Hour).Unix(), // 8 часов жизни токена
	})

	// Подписываем токен секретным ключом
	return token.SignedString(secretKey)
}

// ValidateJWT проверяет токен и извлекает хэш пароля
func ValidateJWT(tokenString string) (string, error) {
	secretKey := getSecretKey()
	if len(secretKey) == 0 {
		return "", errors.New("TODO_PASSWORD is not set")
	}

	// Разбираем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}

	// Читаем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if hash, ok := claims["password_hash"].(string); ok {
			return hash, nil
		}
	}

	return "", errors.New("invalid token")
}
