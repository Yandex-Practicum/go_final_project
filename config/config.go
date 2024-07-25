package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Сonfig struct {
	Port   string
	DBFile string
}

func New() (*Сonfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	cfg := Сonfig{
		Port:   os.Getenv("TODO_PORT"),
		DBFile: os.Getenv("TODO_DBFILE"),
	}

	return &cfg, nil
}
