package handler

import (
	"log"
	"net/http"
)

const webDir = "./web"

func GetFront() http.Handler {
	log.Printf("Loaded frontend files from %s\n", webDir)

	return http.FileServer(http.Dir(webDir))
}
