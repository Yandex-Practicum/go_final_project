package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Сonfig struct {
	Port     string
	DBFile   string
	Password string
	LogLevel log.Level
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
	}

	logLevel, err := log.ParseLevel(os.Getenv("TODO_LOGLEVEL"))
	if err != nil {
		return nil, err
	}

	cfg.LogLevel = logLevel

	return &cfg, nil
}
