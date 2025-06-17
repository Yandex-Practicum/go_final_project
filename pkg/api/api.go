package api

import (
	"net/http"
)

func Init() {

	http.HandleFunc("/api/nextdate", NextDateHandler)

	http.HandleFunc("/api/task", taskHandler)

}
