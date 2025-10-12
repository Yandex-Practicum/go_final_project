package api

import (
	"os"

	"go_final_project/pkg/auth"
)

// Config структура для хранения конфигурации
type Config struct {
	Password string
	Port     string
	DBfile   string
}

// appConfig глобальная переменная конфигурации
var appConfig Config

// InitConfig инициализирует конфигурацию при старте приложения
func InitConfig() {
	appConfig = Config{
		Password: os.Getenv("TODO_PASSWORD"),
		Port:     getEnvWithDefault("TODO_PORT", "7540"),
		DBfile:   getEnvWithDefault("TODO_DBFILE", "scheduler.db"),
	}

	auth.Init(appConfig.Password)

}

// GetConfig возвращает конфигурацию (read-only)
func GetConfig() Config {
	return appConfig
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
