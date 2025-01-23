package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// LoadEnv загружает переменные окружения из .env файла.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err) // Логируем и завершаем программу
	}
}

// MakeDB создаёт или открывает базу данных
func MakeDB() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Ошибка получения пути к исполняемому файлу: %v", err)
	}

	// Получаем путь к файлу базы данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db") // Дефолтный путь
		log.Printf("Путь к базе данных не задан. Используется путь по умолчанию: %s", dbFile)
	}

	// Проверяем, существует ли файл базы данных
	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
		log.Printf("Файл базы данных не найден. Будет создан новый: %s", dbFile)
	}

	// Открываем базу данных
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Если база данных новая, создаем схему
	if install {
		err = createSchema(DB)
		if err != nil {
			log.Fatalf("Ошибка создания схемы базы данных: %v", err)
		}
		log.Println("База данных создана")
	} else {
		log.Println("База данных уже существует")
	}
}

// createSchema создаёт таблицы и индекс в базе данных
func createSchema(db *sql.DB) error {
	// Создаем таблицу и индекс
	query := `
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT CHECK(length(repeat) <= 128)
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
    `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("Схема базы данных создана или уже существует")
	return nil
}

// CloseDB закрывает соединение с базой данных.
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
		log.Println("База данных закрыта")
	}
}
