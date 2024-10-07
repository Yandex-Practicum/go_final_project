package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var jwtKey = []byte("my_secret_key")

func main() {
	// Получаем порт из переменной окружения TODO_PORT
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // Устанавливаем порт по умолчанию, если переменная не задана
	}

	// Определяем путь к базе данных
	dbPath := getDatabasePath()

	// Подключаемся к базе данных
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализируем базу данных
	if err := initializeDatabase(db); err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}

	// Директория с файлами фронтенда
	webDir := "./web"

	// Настраиваем файловый сервер для обслуживания статических файлов
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// Настраиваем API-обработчики с аутентификацией
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", authMiddleware(makeHandler(taskHandler, db)))
	http.HandleFunc("/api/tasks", authMiddleware(makeHandler(tasksHandler, db)))
	http.HandleFunc("/api/task/done", authMiddleware(makeHandler(taskDoneHandler, db)))
	http.HandleFunc("/api/signin", makeHandler(signInHandler, db))

	// Запуск сервера на указанном порту
	log.Printf("Запуск сервера на порту %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// getDatabasePath возвращает путь к файлу базы данных
func getDatabasePath() string {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbPath = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}
	log.Println("Путь к базе данных:", dbPath)
	return dbPath
}

// initializeDatabase инициализирует базу данных и создает таблицы, если необходимо
func initializeDatabase(db *sql.DB) error {
	// Проверяем существование файла базы данных
	_, err := os.Stat(getDatabasePath())
	install := false
	if os.IsNotExist(err) {
		install = true
	}

	// Если базы данных нет, создаем таблицу scheduler
	if install {
		createTableQuery := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT(128)
		);
		CREATE INDEX idx_date ON scheduler(date);
		`
		_, err := db.Exec(createTableQuery)
		if err != nil {
			return err
		}
		log.Println("База данных и таблица созданы.")
	}

	return nil
}

// makeHandler оборачивает обработчик и передает ему базу данных
func makeHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}
