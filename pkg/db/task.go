package db

import (
	"database/sql"
	"fmt"
	"log"
)

// подключение к бд для инициализации в main.go
func Init(database *sql.DB) {
	db = database
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func GetTask(id string) (*Task, error) {
	log.Println("GetTask", id)
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	var task Task

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если задача не найдена, возвращаем ошибку
			return nil, fmt.Errorf("task not found")
		}
		// Если ошибка не связана с отсутствием строки, возвращаем ошибку
		return nil, err
	}

	return &task, nil
}

// Функция для поиска задач с поддержкой поиска по заголовку, комментарию и дате
func Tasks(limit int, search string, dateSearch string) ([]*Task, error) {
	// Базовый запрос
	query := `SELECT id, date, title, comment, repeat FROM scheduler`

	// Если есть параметр search, добавляем условие LIKE
	if search != "" {
		query += ` WHERE title LIKE ? OR comment LIKE ?`
	}

	// Если есть параметр dateSearch (поиск по дате), добавляем условие для даты
	if dateSearch != "" {
		if search != "" {
			query += ` AND date = ?`
		} else {
			query += ` WHERE date = ?`
		}
	}

	// Добавляем сортировку по дате
	query += ` ORDER BY date LIMIT ?`

	// Подготовка параметров для запроса
	var params []interface{}
	if search != "" {
		params = append(params, "%"+search+"%", "%"+search+"%")
	}
	if dateSearch != "" {
		params = append(params, dateSearch)
	}
	params = append(params, limit)

	// Выполнение запроса
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	// Если задач нет, возвращаем пустой слайс
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}
