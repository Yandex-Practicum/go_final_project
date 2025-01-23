package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

// SignUpHandler обрабатывает запросы на регистрацию
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// Определяем структуру для хранения данных из запроса
	var payload struct {
		Password string `json:"password"`
	}
	// Декодируем JSON-данные из тела запроса в структуру payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"Некорректный запрос"}`, http.StatusBadRequest)
		return
	}
	// Получаем ожидаемый пароль из переменных окружения
	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" || payload.Password != envPassword {
		http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
		return
	}
	// Создаем новый JWT с указанием метода подписи и утверждений (claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password":   envPassword,
		"expiration": time.Now().Add(time.Hour * 8).Unix(),
	})
	// Подписываем токен с использованием секретного ключа
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, `{"error": "Ошибка генерации токена"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": signedToken,
	})
}

// AuthMiddleware прослойка для аутентификации
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		envPassword := os.Getenv("TODO_PASSWORD")
		if envPassword == "" {
			next.ServeHTTP(w, r)
			return
		}
		// Получаем токен из cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, `{"error": "Необходима аутентификация"}`, http.StatusUnauthorized)
			return
		}
		// Парсим токен и проверяем его подпись
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Токен не действителен"}`, http.StatusUnauthorized)
			return
		}
		// Извлекаем утверждения (claims) из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["password"] != envPassword {
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
