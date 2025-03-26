package main

import (
	"log"

	"todo/config"
	"todo/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
