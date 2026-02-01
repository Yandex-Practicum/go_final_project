package config

import (
	"context"
	"errors"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	WebDir   string `env:"TODO_WEB_DIR,default=./web"`
	Port     string `env:"TODO_PORT,default=:7540"`
	DbFile   string `env:"TODO_DBFILE,default=./pkg/db/scheduler.db"`
	Password string `env:"TODO_PASSWORD,default=12345"`
}

func NewConfig() (*Config, error) {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(".env"); err != nil {
		err = errors.New("warning: .env file not found or couldn't be loaded")
		return nil, err
	}

	cfg := &Config{}

	// Создаем пустой контекст, т.к. этого требует парсер envconfig.Process
	ctx := context.Background()

	// Парсим переменные окружения в структуру Config
	if err := envconfig.Process(ctx, cfg); err != nil {
		err = errors.New("failed to parse env config: " + err.Error())
		return nil, err
	}

	return cfg, nil
}
