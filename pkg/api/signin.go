package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("failed to decode JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJson(w, map[string]string{"error": "authentication not configured"}, http.StatusBadRequest)
		return
	}

	if input.Password != expectedPassword {
		writeJson(w, map[string]string{"error": "invalid password"}, http.StatusUnauthorized)
		return
	}

	hash := sha256.Sum256([]byte(expectedPassword))
	hashStr := hex.EncodeToString(hash[:])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password_hash": hashStr,
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(expectedPassword))
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("failed to create token: %v", err)}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{"token": tokenString}, http.StatusOK)
}
