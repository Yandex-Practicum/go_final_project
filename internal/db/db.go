package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// SetupDB инициализирует базу данных и создает таблицу, если она не существует
func SetupDB() (*sqlx.DB, error) {
	// Получаем путь из переменной окружения
	dbPath := os.Getenv("TODO_DBFILE")

	// Если переменная окружения не указана
	if dbPath == "" {
		// Получаем текущую директорию проекта
		projectDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения рабочей директории: %w", err)
		}
		dbPath = filepath.Join(projectDir, "scheduler.db")                     // Создаем путь в директории проекта
		log.Println("Используется путь по умолчанию для базы данных:", dbPath) // Логируем использование пути по умолчанию
	} else {
		log.Println("Используется путь к базе данных из переменной окружения:", dbPath) // Логируем использование переменной окружения
	}

	// Проверяем, существует ли файл базы данных
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Если файла нет, создаем базу данных и таблицу
		db, err := sqlx.Connect("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
		}

		// Создаем таблицу и индекс
		createTableQuery := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128)
		);
		CREATE INDEX idx_date ON scheduler(date);
		`
		if _, err := db.Exec(createTableQuery); err != nil {
			db.Close()
			return nil, fmt.Errorf("ошибка создания таблицы: %v", err)
		}

		log.Printf("Создана база данных по пути: %s\n", dbPath) // Логируем создание базы данных
		return db, nil
	}

	// Если файл уже существует, просто открываем базу данных
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	log.Printf("Открыта база данных по пути: %s\n", dbPath) // Логируем открытие базы данных
	return db, nil
}
