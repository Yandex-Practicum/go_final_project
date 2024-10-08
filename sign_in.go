package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Обработчик для /api/signin
func signInHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	var creds struct {
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, `{"error":"Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}
	if appPassword == "" || creds.Password != appPassword {
		http.Error(w, `{"error":"Неверный пароль"}`, http.StatusUnauthorized)
		return
	}
	tokenString, err := generateToken()
	if err != nil {
		http.Error(w, `{"error":"Ошибка генерации токена"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
