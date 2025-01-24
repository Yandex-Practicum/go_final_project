package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sergeysmagin95/go_final_project/tests"
)

func main() {
	logCreateDatabase := CreateDatabase()
	log.Println(logCreateDatabase)

	// Инициализация базы данных
	if err := initDB(); err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close() // Закрываем соединение при завершении работы приложения

	http.HandleFunc("/api/nextdate", nextDateHandler)  // Обработчик для /api/nextdate
	http.HandleFunc("/api/tasks", tasksHandler)        // Обработчик для /api/tasks
	http.HandleFunc("/api/task", taskHandler)          // Обработчик для добавления, получения, обновления, удаления задачи
	http.HandleFunc("/api/task/done", doneTaskHandler) // Обработчик для пометки задачи как выполненной

	//Устанавливаем порт по умолчанию
	port := strconv.Itoa(tests.Port)

	//Проверяем переменную окруженния TODO_PORT
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	//Устанавливаем директорию для файлов
	webDir := "./web"

	//Настраиваем обработчик для статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	//Запускаем сервер
	log.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// CreateDatabase создает базу данных и таблицу scheduler
func CreateDatabase() string {
	// Получаем путь к файлу базы данных из переменной окружения
	dbFileName := os.Getenv("TODO_DBFILE")
	if dbFileName == "" {
		dbFileName = "scheduler.db" // Путь по умолчанию
	}
	// Проверяем, существует ли файл базы данных
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// SQL-запрос для создания таблицы
		createTableSQL := `
		BEGIN;
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat VARCHAR(128)
        );
		CREATE INDEX date_idx ON scheduler (date);
		COMMIT;`

		// Выполняем запрос
		if _, err := db.Exec(createTableSQL); err != nil {
			log.Fatalf("Ошибка создания таблицы: %v", err)
		}

		return "База данных и таблица успешно созданы."
	} else {
		return "База данных уже существует."
	}
}
