package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte

// Init инициализирует JWT секрет
func Init(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
}

// Claims структура для JWT токена
type Claims struct {
	PasswordHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

// GenerateToken создает JWT токен
func GenerateToken() (string, error) {
	if len(jwtSecret) == 0 {
		return "", fmt.Errorf("JWT secret not configured")
	}

	// Создаем claims с хешем пароля
	claims := &Claims{
		PasswordHash: fmt.Sprintf("%x", jwtSecret), // для демонстрации
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}

// ValidateToken проверяет JWT токен
func ValidateToken(tokenString string) (bool, error) {
	if len(jwtSecret) == 0 {
		return true, nil // Аутинтификация отключена
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Проверяем что хеш пароля совпадает с текущим
		expectedHash := fmt.Sprintf("%x", jwtSecret)
		return claims.PasswordHash == expectedHash, nil
	}

	return false, fmt.Errorf("invalid token")

}

// IsAuthEnabled проверяет включена ли аутентификация
func IsAuthEnabled() bool {
	return len(jwtSecret) > 0
}
