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

// GenerateToken creates a new JWT token with an 8-hour expiration time.
func GenerateToken() (string, error) {
	if len(secretKey) == 0 {
		return "", errors.New("secret key is empty")
	}

	// Creating a token with a lifetime of 8 hours
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	// Token generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken checks if the given JWT token is valid.
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

// AuthMiddleware ensures authentication before granting API access.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if len(secretKey) == 0 {
			next(w, req)
			return
		}
		cookie, err := req.Cookie("token")
		if err != nil {
			sendAuthError(w, "token is invalid")
			return
		}
		if err := ValidateToken(cookie.Value); err != nil {
			sendAuthError(w, "token is invalid")
			return
		}
		next(w, req)
	}
}

// sendAuthError sends a JSON-formatted authentication error response.
func sendAuthError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	jsonResponse := `{"error":"` + message + `"}`
	w.Write([]byte(jsonResponse))
}
