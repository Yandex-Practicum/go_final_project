package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}
	return db, nil
}

func FindOrCreateDB(todoDB string) (string, error) {
	if todoDB == "" {
		todoDB = "scheduler.db" // Файл БД в корне проекта
	}

	if _, err := os.Stat(todoDB); os.IsNotExist(err) {
		if err := createDB(todoDB); err != nil {
			return "", fmt.Errorf("не удалось создать БД: %w", err)
		}
	}

	return todoDB, nil
}

func createDB(dbFile string) error {
	// Создаём файл БД
	file, err := os.Create(dbFile)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	file.Close()

	// Открываем БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}
	defer db.Close()

	// Путь к SQL-скрипту (относительно корня проекта)
	sqlPath := filepath.Join("sqlite", "scheduler_creator.sql")
	textSQL, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения SQL-файла: %w", err)
	}

	// Выполняем скрипт
	if _, err = db.Exec(string(textSQL)); err != nil {
		return fmt.Errorf("ошибка выполнения SQL: %w", err)
	}

	return nil
}
