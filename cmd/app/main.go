package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/JKasus/go_final_project/pkg/api"
	"github.com/JKasus/go_final_project/pkg/config"
	"github.com/JKasus/go_final_project/pkg/db"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	err = db.Init(cfg.DbFile)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	log.Println("database initialized")
	defer db.Close()
	r := chi.NewRouter()
	api.Init(r)
	r.Handle("/*", http.FileServer(http.Dir(cfg.WebDir)))
	log.Println("starting server on %s", cfg.Port)
	if err = http.ListenAndServe(cfg.Port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
