package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON отправляет JSON-ответ
func writeJson(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(value)
}

// writeJSONError отправляет JSON с ошибкой
func writeJsonError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	writeJson(w, map[string]string{"error": message})
}
