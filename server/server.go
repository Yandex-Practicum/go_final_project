package server

import (
	"fmt"
	nextdate "go1f/pkg/api"
	"net/http"
)

func Run() error {
	newr := nextdate.Init()
	port := 7540
	return http.ListenAndServe(fmt.Sprintf(":%d", port), newr)
}
