package common

import (
	"encoding/json"
	"net/http"

	"go_final_project/internal/models"
)

func SetJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func Respond(w http.ResponseWriter, v any) {
	SetJsonHeader(w)
	json.NewEncoder(w).Encode(v)
}

func RespondWithError(w http.ResponseWriter, err error) {
	response := models.Response{Error: err.Error()}
	SetJsonHeader(w)
	json.NewEncoder(w).Encode(response)
}
