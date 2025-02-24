package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

func InsertAndReturnID(db *sql.DB, date, title, comment, repeat string) (string, error) {

	createTaskQuery := `
	INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)`

	log.Printf("Inserting task: Date=%s, Title=%s, Comment=%s, Repeat=%s\n", date, title, comment, repeat)

	// Выполняю запрос с передачей значений переменных
	res, err := db.Exec(createTaskQuery, date, title, comment, repeat)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	log.Printf("id = %d", id)
	return strconv.FormatInt(id, 10), nil
}

// Функция для создания базы данных
func CreateDatabase(dbFile string) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Не удалось создать базу данных: %v", err)
	}
	defer db.Close()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS scheduler (
		"id" INTEGER PRIMARY KEY,
		"date" TEXT NOT NULL DEFAULT "",
		"title" TEXT NOT NULL DEFAULT "",
		"comment" TEXT,
        repeat TEXT CHECK(length(repeat) <= 128)
	);
																						
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	}

	log.Println("База данных и таблица успешно созданы.")
}

// Функция для подключения к бд (и создания бд если ее нет)
func ConnectDB() *sql.DB {

	// Получение значения TODO_DBFILE
	DBFILE := os.Getenv("TODO_DBFILE")
	if DBFILE == "" {
		log.Fatal("Переменная окружения TODO_DBFILE не установлена")
	}

	var install bool

	// проверяю наличие бд и устанавливаю значение для install
	if _, err := os.Stat(DBFILE); err != nil {
		if os.IsNotExist(err) {
			log.Printf("бд нет")
			install = true
		} else {
			log.Fatalf("ошибка: %v", err)
		}
	}
	// создаю бд если install = true
	if install {
		CreateDatabase(DBFILE)
	}
	// Подключение к базе данных
	db, err := sql.Open("sqlite3", DBFILE)
	if err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v", err)
	}

	log.Println("бд подключена")
	return db
}
