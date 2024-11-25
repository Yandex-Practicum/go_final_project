package error

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
