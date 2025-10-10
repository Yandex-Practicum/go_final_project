package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных и возвращает ID
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTasks возвращает список задач с ограничением по количеству
func GetTasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          ORDER BY date ASC, id ASC 
	          LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Если задач нет, возвращаем пустой слайс вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// SearchTasks возвращает задачи по поисковому запросу
func SearchTasks(search string, limit int) ([]*Task, error) {
	fmt.Printf("SearchTasks вызвана с search='%s', limit=%d\n", search, limit)

	// Проверяем, является ли поисковый запрос датой в формате DD.MM.YYYY
	if isDateSearch(search) {
		fmt.Println("Определен поиск по дате")
		date, err := convertSearchDate(search)
		if err != nil {
			fmt.Printf("Ошибка конвертации даты: %v\n", err)
			return nil, err
		}
		fmt.Printf("Конвертированная дата: %s\n", date)
		return searchByDate(date, limit)
	}

	// Обычный поиск по тексту
	fmt.Println("Выполняется текстовый поиск")
	return searchText(search, limit)
}

// isDateSearch проверяет, является ли строка датой в различных форматах
func isDateSearch(s string) bool {
	// Пробуем разные разделители
	for _, separator := range []string{".", ",", "/", "-"} {
		if len(s) == 10 && s[2] == separator[0] && s[5] == separator[0] {
			return true
		}
	}
	return false
}

// convertSearchDate конвертирует дату из различных форматов в YYYYMMDD
func convertSearchDate(dateStr string) (string, error) {
	// Заменяем все возможные разделители на точку для унификации
	normalized := strings.ReplaceAll(dateStr, ",", ".")
	normalized = strings.ReplaceAll(normalized, "/", ".")
	normalized = strings.ReplaceAll(normalized, "-", ".")

	// Парсим дату в формате "02.01.2006"
	t, err := time.Parse("02.01.2006", normalized)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}
	return t.Format("20060102"), nil
}

// searchByDate ищет задачи по дате
func searchByDate(date string, limit int) ([]*Task, error) {
	fmt.Printf("searchByDate: поиск по дате %s\n", date)

	// Исправлено: WHERE вместо WERE и используем простые плейсхолдеры
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          WHERE date = ? 
	          ORDER BY date ASC, id ASC 
	          LIMIT ?`

	rows, err := DB.Query(query, date, limit)
	if err != nil {
		fmt.Printf("Ошибка запроса searchByDate: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

// searchText ищет задачи по тексту в заголовке или комментариях
func searchText(search string, limit int) ([]*Task, error) {
	fmt.Printf("searchText: поиск по тексту '%s'\n", search)

	searchPattern := "%" + search + "%"
	// Используем простые плейсхолдеры
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          WHERE title LIKE ? OR comment LIKE ? 
	          ORDER BY date ASC, id ASC 
	          LIMIT ?`

	rows, err := DB.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		fmt.Printf("Ошибка запроса searchText: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

// scanTasks сканирует строки из результата запроса
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			fmt.Printf("Ошибка сканирования строки: %v\n", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	fmt.Printf("scanTasks: найдено %d задач\n", len(tasks))
	return tasks, nil
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler
	WHERE id = ?`

	task := &Task{}
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}
	return task, nil
}

// UpdateTask обновляет задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil

}

// DeleteTask Удаляет задачу по ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil

}

// UpdateDate обновляет только дату задачи
func UpdateDate(id string, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := DB.Exec(query, newDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil

}
