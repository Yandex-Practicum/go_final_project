package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func CheckCreateDB() (*sql.DB, func(), error) {

	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		return nil, nil, err
	}

	// Определение расположения исполняемого файла
	// Проверка присутствия файла БД
	appPath, err := os.Executable()
	if err != nil {
		return nil, nil, err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), os.Getenv("DB_NAME"))
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true // Установка признака в необходимости создания БД и её инициализации
	}

	// Подключение к БД
	// если БД нет, она создаётся и пингуется
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_NAME"))
	if err != nil {
		return nil, func() { _ = db.Close() }, err
	}

	err = db.Ping()
	if err != nil {
		return nil, func() { _ = db.Close() }, err
	}

	if install {

		str := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
 			"date" VARCHAR(250),
 			"title" VARCHAR(250) NOT NULL DEFAULT '',
 			"comment" VARCHAR(250) NOT NULL DEFAULT '',
 			"repeat" VARCHAR(128)
			);
    	`, os.Getenv("DB_TABLE_NAME"))

		stmt, err := db.Prepare(str)

		if err != nil {
			return nil, func() { _ = db.Close() }, err
		}
		defer func() { _ = stmt.Close() }()

		_, err = stmt.Exec()
		if err != nil {
			return nil, func() { _ = db.Close() }, err
		}

		fmt.Println("Таблица создана.")

		str = fmt.Sprintf("CREATE INDEX %s_date ON %s (date)", os.Getenv("DB_TABLE_NAME"), os.Getenv("DB_TABLE_NAME"))

		_, err = db.Exec(str)
		if err != nil {
			return nil, func() { _ = db.Close() }, err
		}

		fmt.Println("Индекс создан.")
	}

	return db, nil, nil // возврат указателя на созданную БД
}
