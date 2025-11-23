package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const (
	// Константа defaultDBFile задаёт имя файла БД по умолчанию.
	defaultDBFile = "scheduler.db"

	// Константа envDBVar указывает имя переменной окружения, через которую можно задать путь к файлу БД.
	envDBVar = "TODO_DBFILE"
)

// Константы содержат SQL-скрипты для создания таблицы scheduler и индекса по полю date, если они ещё не существуют.
const (
	createTableSQL = `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(255) NOT NULL,
		comment TEXT,
		repeat VARCHAR(128)
	);`
	createIndexSQL = `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);`
)

// Функция Init инициализирует подключение к базе данных SQLite.
// Параметры:
// dbFile - путь к файлу БД (может быть пустым).
// Возвращает:
// *sql.DB - объект подключения к БД,
// error - ошибку, если инициализация не удалась.
// Логика работы:
//  1. Определяет путь к БД: сначала проверяет переданный аргумент, затем переменную окружения TODO_DBFILE, затем использует значение по умолчанию.
//  2. Проверяет существование файла БД.
//  3. Открывает соединение с БД и настраивает параметры подключения.
//  4. Проверяет доступность БД (ping).
//  5. Если БД не существовала - создаёт схему (таблицу и индекс).
func Init(dbFile string) (*sql.DB, error) {
	// Определяем путь к БД: приоритет - переданный аргумент, затем env, затем дефолт
	if dbFile == "" {
		dbFile = os.Getenv(envDBVar)
	}
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	// Проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return nil, fmt.Errorf("failed to access database file %q: %w", dbFile, err)
		}
	}

	// Открываем соединение с БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка параметров соединения:
	// - максимальное число открытых соединений: 10,
	// - максимальное число idle-соединений: 5,
	// - время жизни соединения: 30 минут.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Проверяем подключение к БД (выполняем пинг к серверу)
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Если БД не существовала - создаём схему (таблицу и индекс)
	if install {
		// Выполняем SQL-скрипт создания схемы
		if _, err = db.Exec(createTableSQL); err != nil {
			// Закрываем соединение при ошибке создания схемы
			db.Close()
			// Отдельная ошибка для CREATE TABLE
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
		// Выполняем SQL-скрипт создания индекса по полю date
		if _, err = db.Exec(createIndexSQL); err != nil {
			// Закрываем соединение при ошибке создания индекса
			db.Close()
			// Отдельная ошибка для CREATE INDEX
			return nil, fmt.Errorf("failed to create index: %w", err)
		}
		log.Println("База данных инициализирована: таблица и индекс созданы")
	} else {
		log.Println("База данных уже существует, схема проверена")
	}

	// Возвращаем готовое соединение с БД
	return db, nil
}
