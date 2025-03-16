package configs

import (
	"os"
)

type Server struct {
	Port string
}

type Database struct {
	DriverName   string
	DatabaseName string
}

type Config struct {
	Server   *Server
	Database *Database
}

func New() *Config {
	return &Config{
		Server: &Server{
			Port: getEnv("TODO_PORT", "7540"),
		},
		Database: &Database{
			DatabaseName: getEnv("TODO_DBFILE", "scheduler.db"),
			DriverName:   getEnv("TODO_DBDRIVER", "sqlite3"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	return defaultVal
}
