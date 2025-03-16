package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go_final_project-main/internal/config"

	"github.com/jmoiron/sqlx"
)

func CheckDb(db *sqlx.DB, config *config.Config) error {
	// Создаем путь к файлу базы данных
	dbFile := config.DbFile

	// Если путь относительный, преобразуем его в абсолютный
	if !filepath.IsAbs(dbFile) {
		appPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("error checking db for existence: %+v", err)
		}
		dbFile = filepath.Join(filepath.Dir(appPath), dbFile)
	}

	// Проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)

	var install bool
	if os.IsNotExist(err) {
		install = true
	}

	// Если база данных новая, создаем таблицу и индекс
	if install {
		schema := `
			CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date TEXT NOT NULL,
				title TEXT NOT NULL,
				comment TEXT,
				repeat TEXT CHECK(LENGTH(repeat) <= 128)
			);
			CREATE INDEX idx_date ON scheduler(date);
		`
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("error creating index: %+v", err)
		}
		log.Printf("Создана новая база данных и таблица %s", dbFile)
	} else {
		log.Printf("Подключено к существующей базе данных %s", dbFile)
	}

	return nil
}
