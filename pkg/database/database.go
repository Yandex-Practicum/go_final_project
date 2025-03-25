package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const dbFileStart = "scheduler.db"

// getDBPath возвращает путь к файлу базы данных.
// Если задана переменная окружения TODO_DBFILE, используем её.
// Иначе возвращаем значение по умолчанию "scheduler.db".
func getDBPath() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = dbFileStart
	}
	return dbFile
}

// InitDB проверяет наличие фалйла DB и создает его при отстувии
func InitDB() (*sql.DB, error) {
	//Проверяем существование файла базы данных
	dbFile := getDBPath()
	//if _, err := os.Stat(dbFile); os.IsNotExist(err) {
	//	log.Println("Файл базы дпнных не найден. Создаем новый файл...")
	//
	//	//Сощдаем фафлй базы данных
	//	file, err := os.Create(dbFile)
	//	if err != nil {
	//		return nil, fmt.Errorf("не удалось создать файл базы данных: %v", err)
	//	}
	//	err = file.Close()
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	//Открываем базу данных
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу данных: %v", err)
	}
	//defer func(db *sql.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		log.Printf("не удалось закрыть базу данных: %v", err)
	//	}
	//}(db)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе: %v", err)
	}

	//Создаем таблицу scheduler,если ее нет
	query := `
			CREATE TABLE IF NOT EXISTS scheduler (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		date TEXT NOT NULL,
    		title TEXT NOT NULL,
   			 comment TEXT NOT NULL,
   			 repeat TEXT NOT NULL
									);
		CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
			`

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать таблицу: %v", err)

	}

	log.Println("База данных успешно инициализирована")

	return db, nil

}
