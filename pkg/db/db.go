package db

import (
	"database/sql"
	"fmt"
	"os"

	// импорт драйвера без использования Cgo

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// schema: создание таблицы и индекса, если их нет
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8)   NOT NULL DEFAULT '',
    title   VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat  VARCHAR(128) DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// Init открывает или создаёт базу данных по пути dbFile.
// Если файла не было, выполняет схему schema.
func Init(dbFile string) error {
	// проверяем существование файла
	_, err := os.Stat(dbFile)     // os.Stat для проверки файла
	install := os.IsNotExist(err) // true, если файл отсутствует :contentReference[oaicite:5]{index=5}

	// открываем БД через database/sql + modernc.org/sqlite
	db, err := sql.Open("sqlite", dbFile) // driver "sqlite" зарегистрирован импортом :contentReference[oaicite:6]{index=6}
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if install {
		// создаём таблицу и индекс
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to init schema: %w", err)
		}
	}

	DB = db
	return nil
}

func DeleteTask(id string) error {
	var idInt int64
	if _, err := fmt.Sscan(id, &idInt); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, idInt)
	if err != nil {
		return fmt.Errorf("db delete error: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next, id string) error {
	var idInt int64
	if _, err := fmt.Sscan(id, &idInt); err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	res, err := DB.Exec(
		`UPDATE scheduler SET date = ? WHERE id = ?`,
		next, idInt,
	)
	if err != nil {
		return fmt.Errorf("db update date error: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
