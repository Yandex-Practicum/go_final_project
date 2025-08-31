package main

import (
	"go_final_project/pkg/api"
	"log"
	"net/http"

	"go_final_project/pkg/db"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("db close: %v", err)
		}
	}()

	api.Init()

	const webDir = "web"
	const addr = ":7540"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("DB: ./scheduler.db ready, serving %s on http://localhost%s/", webDir, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
