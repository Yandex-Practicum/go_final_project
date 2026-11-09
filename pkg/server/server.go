package server

import (
	"database/sql"
	"github.com/Yandex-Practicum/go_final_project/pkg/api"
	"log"
	"net/http"
)

func Start(database *sql.DB) {
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init(database)

	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatal(err)
	}

}
