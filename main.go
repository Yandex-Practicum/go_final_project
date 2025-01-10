package main

import (
	//"fmt"
	"net/http"
	//"os"
	//"github.com/go-chi/chi/v5"
)

var (
	port   = ":7540"
	webDir = "./web"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}

}
