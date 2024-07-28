package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Сonfig struct {
	Port     string
	DBFile   string
	Password string
	LogLevel string
}

func New() (*Сonfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	cfg := Сonfig{
		Port:     os.Getenv("TODO_PORT"),
		DBFile:   os.Getenv("TODO_DBFILE"),
		Password: os.Getenv("TODO_PASSWORD"),
		LogLevel: os.Getenv("TODO_LOGLEVEL"),
	}

	return &cfg, nil
}
