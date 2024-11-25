package database

import (
	"log"
	"os"
	"path/filepath"

	"go_final_project/internal/repository"
)

func Migration(rep *repository.Repository) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	if install {
		query := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(256) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT "");
			CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`

		if _, err := New().Exec(query); err != nil {
			log.Println("Ошибка создания таблицы")
			return
		}
	}
	log.Println("Таблица успешно создана")
}
