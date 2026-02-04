package config

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

type Config struct {
	Port      string `arg:"env:TODO_PORT" default:"7540"`
	DBFile    string `arg:"env:TODO_DBFILE" default:"scheduler.db"`
	Password  string `arg:"env:TODO_PASSWORD" default:"password"`
	JWTSecret string `arg:"env:TODO_JWT_SECRET" default:"secret"`
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := cfg.loadFromEnv(); err != nil {
		return nil, fmt.Errorf("unable to load config")
	}

	return cfg, nil
}

func (c *Config) loadFromEnv() error {
	return arg.Parse(c)
}
