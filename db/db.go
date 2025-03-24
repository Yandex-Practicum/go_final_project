package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Импорт драйвера SQLite
)

// Task представляет задачу в базе данных
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var db *sql.DB // Объявляем переменную для базы данных

// InitDB инициализирует базу данных и создает необходимые таблицы
func InitDB(isTest bool) (*sql.DB, error) {
	var dbFile string

	if isTest {
		dbFile = filepath.Join(os.TempDir(), "test_scheduler.db")
	} else {
		dbFile = "./scheduler.db"
	}

	if isTest {
		if _, err := os.Stat(dbFile); err == nil {
			os.Remove(dbFile)
		}
	}

	var err error
	db, err = sql.Open("sqlite", dbFile) // Присваиваем значение переменной db
	if err != nil {
		return nil, err
	}

	var tableExists string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler';").Scan(&tableExists)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Создаем таблицу, если она не существует
	createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	// Создаем индекс
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetDB возвращает текущее соединение с базой данных
func GetDB() (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}
	return db, nil
}

func GetTaskByID(id string) (Task, error) {
	var task Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("задача не найдена")
		}
		log.Printf("Ошибка при выполнении запроса: %v", err) // Логируем ошибку
		return Task{}, err
	}
	return task, nil
}

func UpdateTask(task Task) error {
	// Проверяем, существует ли задача с указанным идентификатором
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM scheduler WHERE id = ?)", task.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования задачи: %w", err)
	}
	if !exists {
		return fmt.Errorf("задача с ID %s не найдена", task.ID)
	}

	// Выполняем обновление задачи
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}
	return nil
}

func DeleteTaskByID(id int) error {
	db, err := GetDB() // Получаем соединение с базой данных
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}
