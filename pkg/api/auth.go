package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

var jwtSecret = []byte("some-secret-key") // пока можно захардкодить

func someHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func generateJWT(pass string) (string, error) {
	claims := jwt.MapClaims{
		"hash": someHash(pass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJson(w, SignInResponse{Error: "Некорректный запрос"}, http.StatusBadRequest)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJson(w, SignInResponse{Error: "Аутентификация отключена"}, http.StatusBadRequest)
		return
	}

	if req.Password != pass {
		writeJson(w, SignInResponse{Error: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(pass)
	if err != nil {
		writeJson(w, SignInResponse{Error: "Ошибка генерации токена"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, SignInResponse{Token: token}, http.StatusOK)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r) // если пароль не задан — пропускаем без проверки
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		hashInToken, ok := claims["hash"].(string)
		if !ok {
			http.Error(w, "Invalid token data", http.StatusUnauthorized)
			return
		}

		if hashInToken != someHash(pass) {
			http.Error(w, "Invalid token hash", http.StatusUnauthorized)
			return
		}

		next(w, r) // всё хорошо — пропускаем к следующему обработчику
	})
}
