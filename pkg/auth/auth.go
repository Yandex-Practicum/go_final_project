package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

const jwtSecret = "my-secret-key" // Секретный ключ для подписи токена

// SignInRequest структура для запроса на вход
type SignInRequest struct {
	Password string `json:"password"`
}

// SignInResponse структура для ответа
type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// HashPassword создаёт хэш пароля
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// GenerateJWT генерирует JWT-токен
func GenerateJWT(passwordHash string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passwordHash": passwordHash,
		"exp":          time.Now().Add(8 * time.Hour).Unix(), // Время жизни токена 8 часов
	})
	return token.SignedString([]byte(jwtSecret))
}

// ValidateJWT проверяет JWT-токен
func ValidateJWT(tokenString, expectedPassword string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if passwordHash, ok := claims["passwordHash"].(string); ok {
			return passwordHash == HashPassword(expectedPassword), nil
		}
	}
	return false, nil
}

// SignInHandler обрабатывает запрос на вход
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body")
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		respondWithError(w, "Password not configured")
		return
	}

	if req.Password != expectedPassword {
		respondWithError(w, "Неверный пароль")
		return
	}

	token, err := GenerateJWT(HashPassword(req.Password))
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		respondWithError(w, "Internal server error")
		return
	}

	resp := SignInResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func respondWithError(w http.ResponseWriter, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(SignInResponse{Error: errorMsg})
}
