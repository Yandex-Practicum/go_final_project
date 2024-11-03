package main

import (
	"net/http"
)

func main() {

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}

}
