package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var request struct {
	Password string `json:"password"`
}

var tokenResponse struct {
	Token string `json:"token"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, `{"error":"can't read body"}`, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		http.Error(w, `{"error":"can't unmarshal body"}`, http.StatusBadRequest)
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")
	if storedPassword == "" {
		http.Error(w, `{"error":"password is empty"}`, http.StatusInternalServerError)
		return
	}

	if request.Password != storedPassword {
		http.Error(w, `{"error":"invalid password"}`, http.StatusUnauthorized)
		return
	}

	tokenResponse.Token, err = generateTokenJWT(storedPassword)
	if err != nil {
		http.Error(w, `{"error":"can't generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(tokenResponse); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
}

func generateTokenJWT(password string) (string, error) {
	claims := jwt.MapClaims {
		"password_hash": fmt.Sprintf("%x", hashPassword(password)),
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("TODO_PASSWORD"))
	return token.SignedString(secretKey)
}

func hashPassword(password string) []byte {
	return []byte(password)
}

func isValidTokenJWT(token, password string) bool {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(password), nil
	})
	if err != nil || !parsedToken.Valid {
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	storedHash := claims["password_hash"].(string)
	return storedHash == fmt.Sprintf("%x", hashPassword(password))
}

func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			var tokenJWT string

			cookie, err := r.Cookie("token")
			if err != nil {
				tokenJWT = cookie.Value
			}

			if tokenJWT == "" || !isValidTokenJWT(tokenJWT, password) {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
