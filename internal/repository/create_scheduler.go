package repository

import (
	"fmt"
)

// Функция для создания таблицы в базе данных для хранения информации о задачах
func (r *Repository) CreateScheduler() error {
	query := `CREATE TABLE IF NOT EXISTS scheduler (
	   id INTEGER PRIMARY KEY AUTOINCREMENT,        
	   date CHAR(8) NOT NULL DEFAULT "",
	   title VARCHAR(256) NOT NULL DEFAULT "",
	   comment TEXT NOT NULL DEFAULT "",
	   repeat VARCHAR(128) NOT NULL DEFAULT "");`

	if _, err := r.db.Exec(query); err != nil {
		return err
	}
	fmt.Println("База данных и таблица созданы.")

	return nil
}

//id INTEGER PRIMARY KEY AUTOINCREMENT - Целочисленный первичный ключ с автоинкрементом
//date CHAR(8) - Дата в формате YYYYMMDD (строка фиксированной длины 8 символов)
//title VARCHAR(256) - Заголовок задачи (строка переменной длины до 256 символов)
//comment TEXT - Комментарий к задаче (текстовое поле)
//repeat VARCHAR(128) - Правило повторения задачи (строка переменной длины до 128 символов)
