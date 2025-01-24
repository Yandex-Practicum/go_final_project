package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	ID       int64  `json:"id,omitempty"`
	NextDate string `json:"next_date,omitempty"`
	Tasks    []Task `json:"tasks,omitempty"`
	Error    string `json:"error,omitempty"`
}

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Config struct {
	WebDir string
	Port   string
	DBPath string
}

// openDB открывает базу данных и проверяет подключение
func OpenDB(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу данных: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	return db, nil
}

// tableExists проверяет, существует ли таблица в базе данных
func TableExists(db *sql.DB, tableName string) (bool, error) {
	row := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?;`, tableName)
	var name string
	err := row.Scan(&name)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("ошибка при проверке существования таблицы: %w", err)
	}
	return strings.EqualFold(name, tableName), nil
}

// createTable создает таблицу в базе данных
func CreateTable(db *sql.DB) error {
	tableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT CHECK(LENGTH(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`
	_, err := db.Exec(tableSQL)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}
	log.Println("Таблица 'scheduler' создана.")
	return nil
}

// withDB оборачивает обработчик, предоставляя ему доступ к базе данных
func WithDB(handler func(db *sql.DB, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := OpenDB("./scheduler.db")
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "Ошибка базы данных")
			return
		}
		defer db.Close()
		handler(db, w, r)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Ошибка при кодировании ответа в JSON", http.StatusInternalServerError)
	}
}

// writeError отправляет сообщение об ошибке в формате JSON
func WriteError(w http.ResponseWriter, status int, errMsg string) {
	WriteJSON(w, status, map[string]string{"error": errMsg})
}
