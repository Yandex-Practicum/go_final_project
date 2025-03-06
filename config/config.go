package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetConfig() Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		panic("SERVER_PORT is not set")
	}
	return Config{
		ServerPort: port,
	}
}
