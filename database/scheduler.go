package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"

	"TODOGo/config"
	"TODOGo/schemas"
)

// AddTaskToDB вставляет новую задачу в базу данных.
func AddTaskToDB(task schemas.Table) (uint64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := config.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка вставки задачи: %v", err)
	}

	// Получаем ID последней вставленной задачи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID последней вставленной задачи: %v", err)
	}

	log.Printf("Задача добавлена с ID: %d", id)
	return uint64(id), nil
}

// SearchTasksInDB ищет задачи в базе данных.
func SearchTasksInDB(search string, limit int) ([]schemas.Table, error) {
	var query string
	var args []interface{}

	if search == "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
		args = append(args, limit)
	} else {
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?"
			args = append(args, parsedDate.Format("20060102"), limit)
		} else {
			likePattern := "%" + search + "%"
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			args = append(args, likePattern, likePattern, limit)
		}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса к базе данных: %v", err)
	}
	defer rows.Close()

	var tasks []schemas.Table
	for rows.Next() {
		var task schemas.Table
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %v", err)
	}
	return tasks, nil
}

// FetchTaskByID извлекает задачу по идентификатору из базы данных.
func FetchTaskByID(id int) (schemas.Table, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?"
	row := config.DB.QueryRow(query, id)

	var task schemas.Table
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return schemas.Table{}, errors.New("задача не найдена")
		}
		return schemas.Table{}, errors.New("ошибка при извлечении задачи")
	}
	return task, nil
}

// UpdateTask обновляет функцию
func UpdateTask(task schemas.Table) error {
	query := `
	UPDATE scheduler 
	SET date=?, title=?, comment=?, repeat=? 
	WHERE id=?`
	res, err := config.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteTaskByID удаляет задачу из базы данных по идентификатору.
func DeleteTaskByID(id int) error {
	query := "DELETE FROM scheduler WHERE id=?"
	_, err := config.DB.Exec(query, id)
	return err
}

// UpdateTaskDate обновляет дату задачи с указанным идентификатором.
func UpdateTaskDate(id uint64, newDate string) error {
	query := "UPDATE scheduler SET date=? WHERE id=?"
	_, err := config.DB.Exec(query, newDate, id)
	return err
}
