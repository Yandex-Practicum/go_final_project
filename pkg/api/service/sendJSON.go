package service

import (
	"encoding/json"
	"net/http"
)

func SendOkJSONBytes(w http.ResponseWriter, status int, b []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(b)
}

func SendErrorJSON(w http.ResponseWriter, errorString string, statusCode int) {
	type response struct {
		Error string `json:"error"`
	}
	var res response
	res.Error = errorString

	jsByte, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Ошибка подготовки ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	w.Write(jsByte)

}
