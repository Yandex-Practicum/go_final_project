package controller

import (
	"encoding/json"
	"net/http"
)

func ResponseError(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": s})
}
