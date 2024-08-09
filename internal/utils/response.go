package utils

import (
	"encoding/json"
	"net/http"

	"go_final_project/internal/models"
)

func RespondWithError(w http.ResponseWriter, message string) {
	response := models.Response{Error: &message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}
