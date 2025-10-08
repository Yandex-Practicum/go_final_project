package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

// Config содержит конфигурационные параметры
type Config struct {
	Port   string `env:"TODO_PORT" envDefault:"7540"`
	DbFile string `env:"TODO_DBFILE" envDefault:"scheduler.db"`
}

func Read() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error with reading config file: %+v", err)
	}

	return cfg, nil
}
