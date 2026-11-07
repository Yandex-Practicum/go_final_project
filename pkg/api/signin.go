package api

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
)

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Метод не поддерживается",
		})
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Некорректный запрос",
		})
		return
	}

	password := configuredPassword

	if password == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Пароль не настроен на сервере",
		})
		return
	}

	if subtle.ConstantTimeCompare(
		[]byte(request.Password),
		[]byte(password),
	) != 1 {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Неверный пароль",
		})
		return
	}

	token, err := createToken(password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Не удалось создать токен",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
