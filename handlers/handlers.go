package handlers

import (
	"net/http"
)

func MainHandle(w http.ResponseWriter, r *http.Request) {

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

}
