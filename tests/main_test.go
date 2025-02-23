package tests

import (
	"os"
	"testing"

	"log"

	"github.com/joho/godotenv"
)

// GetTestEnv загружает переменные окружения из файла .env.test
func GetTestEnv() {
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env.test: %v", err)
	}
	log.Println("успешная загрузка  файла .env.test")
}

func TestMain(m *testing.M) {
	// Загрузка переменных окружения перед выполнением тестов
	GetTestEnv()
	log.Println("Загрузка переменных окружения перед выполнением тестов")
	// Выполнение всех тестов в пакете
	code := m.Run()

	// Очистка переменных окружения после выполнения тестов
	os.Clearenv()

	// Возврат кода выполнения тестов
	os.Exit(code)
}
