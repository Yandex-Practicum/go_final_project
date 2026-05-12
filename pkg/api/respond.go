package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type fail struct {
	Error string `json:"error"`
}

func writer(w http.ResponseWriter, resp any) {
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
	}
}

func Succefull(w http.ResponseWriter, resp any) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	writer(w, resp)
}

func Failed(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	writer(w, fail{
		Error: err.Error(),
	})
}
