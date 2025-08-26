package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB инициализирует подключение к базе данных
func InitDB(dbFile string) (*sql.DB, error) {
	// Открываем соединение с БД
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Проверяем, существует ли таблица scheduler
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='scheduler';").Scan(&exists)
	if err != nil || exists == 0 {
		log.Println("⚠️ Таблица scheduler не найдена, создаем её...")
		createDatabase(db)
	}

	log.Println("✅ База данных успешно подключена")
	return db, nil
}

// createDatabase создаёт таблицу, если её нет
func createDatabase(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		title TEXT NOT NULL,
		comment TEXT,
		repeat TEXT
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании таблицы: %v", err)
	}
	log.Println("✅ Таблица scheduler успешно создана (или уже существовала)")
}

// GetDBFile возвращает путь к файлу базы данных в `go_final_project_ref/scheduler.db`
func GetDBFile() string {
	// 1️⃣ Проверяем переменную окружения
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		log.Printf("📂 Используемая БД из переменной окружения: %s", envDBFile)
		return envDBFile
	}

	// 2️⃣ Получаем путь к корневой директории проекта
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatal("❌ Ошибка получения рабочего каталога:", err)
	}

	// Поднимаемся до корня проекта (go_final_project_ref)
	for !isProjectRoot(baseDir) {
		baseDir = filepath.Dir(baseDir)
		if baseDir == "/" {
			log.Fatal("❌ Не удалось найти корень проекта go_final_project_ref")
		}
	}

	// 3️⃣ Формируем путь к базе данных
	dbPath := filepath.Join(baseDir, "scheduler.db")
	log.Printf("📂 Ожидаемый путь к БД: %s", dbPath)

	return dbPath
}

// isProjectRoot проверяет, является ли директория корнем проекта
func isProjectRoot(path string) bool {
	_, err := os.Stat(filepath.Join(path, "go.mod"))
	return err == nil
}
