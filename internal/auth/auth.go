package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func GetSecretKey(key string) {
	secretKey = []byte(key)
}

// GenerateToken создает JWT токен для аутентификации
func GenerateToken() (string, error) {
	// Проверяю, установлен ли пароль
	if len(secretKey) == 0 {
		return "", errors.New("secret key is empty")
	}

	// Создаю токен с временем жизни 8 часов
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	// Генерирую токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

// ValidateToken проверяет валидность токена
func ValidateToken(tokenString string) error {
	if len(secretKey) == 0 {
		return errors.New("secret key is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// AuthMiddleware проверяет JWT перед доступом к API
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if len(secretKey) == 0 {
			// Если пароль не задан, то без авторизации
			next(w, req)
			return
		}
		cookie, err := req.Cookie("token")
		if err != nil {
			http.Error(w, `{"error":"token is invalid"}`, http.StatusUnauthorized)
			return
		}
		if err := ValidateToken(cookie.Value); err != nil {
			http.Error(w, `{"error":"token is invalid"}`, http.StatusUnauthorized)
			return
		}
		next(w, req)
	}
}
