package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JKasus/go_final_project/pkg/entities"
)

func writeJSON(w http.ResponseWriter, statusCode int, msg any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	var resp any

	switch v := msg.(type) {
	case error:
		resp = map[string]interface{}{"error": v.Error()}
	case string:
		resp = map[string]interface{}{"nextDate": v}
	case int64:
		resp = map[string]interface{}{"id": strconv.FormatInt(v, 10)}
	case []entities.Task:
		resp = map[string]interface{}{"tasks": v}
	case *entities.Task:
		resp = v
	case entities.EmptyResponse:
		resp = v
	case *entities.UserData:
		resp = map[string]interface{}{"token": v.Token}
	default:
		resp = map[string]interface{}{"error": "invalid type of message"}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
