package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:7540"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

var defaultConfig = Config{
	Env:         "local",
	StoragePath: "./scheduler.db",
	HTTPServer: HTTPServer{
		Address:     "0.0.0.0:7540",
		Timeout:     4 * time.Second,
		IdleTimeout: 60 * time.Second,
	},
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	var config Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = defaultConfig
	} else {
		err := cleanenv.ReadConfig(configPath, &config)
		if err != nil {
			config = defaultConfig
		}
	}

	return &config
}

func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
