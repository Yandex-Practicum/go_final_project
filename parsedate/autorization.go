package parsedate

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SignRequest struct {
	Password string `json:"password"`
}

type SignResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, head int, message string) {
	w.WriteHeader(head)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func SignHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow_Origin", "http://localhost:7540")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	envPassword := os.Getenv("TODO_PASSWORD")
	if envPassword == "" {
		respondWithError(w, http.StatusInternalServerError, "Auth disabled")
		return
	}

	if req.Password != envPassword {
		respondWithError(w, http.StatusUnauthorized, "Wrong password")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passHash": fmt.Sprintf("%x", sha256.Sum256([]byte(envPassword))),
		"exp":      time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(envPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token error")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(8 * time.Hour),
		Path:    "/",
	})
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
