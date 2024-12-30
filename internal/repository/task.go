package repository

import (
	"database/sql"
	"go_final_project/config"
	"go_final_project/internal/models"
	"log"
	"time"
)

func AddTaskToDB(task *models.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := config.GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTasks(search string) ([]map[string]string, error) {
	db := config.GetDB()
	var rows *sql.Rows
	var err error

	query := "SELECT id, date, title, comment, repeat FROM scheduler"
	var args []interface{}

	// Проверяем параметр search
	if search != "" {
		// Определяем, это поиск по дате или подстроке
		if date, parseErr := parseDate(search); parseErr == nil {
			query += " WHERE date = ?"
			args = append(args, date)
		} else {
			query += " WHERE title LIKE ? OR comment LIKE ?"
			search = "%" + search + "%"
			args = append(args, search, search)
		}
	}

	query += " ORDER BY date ASC LIMIT 50"

	// Выполняем запрос
	rows, err = db.Query(query, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Обрабатываем результаты
	var tasks []map[string]string
	for rows.Next() {
		var id, date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			log.Printf("Row scan error: %v", err)
			return nil, err
		}
		tasks = append(tasks, map[string]string{
			"id":      id,
			"date":    date,
			"title":   title,
			"comment": comment,
			"repeat":  repeat,
		})
	}
	return tasks, nil
}

// Вспомогательная функция для преобразования даты
func parseDate(input string) (string, error) {
	date, err := time.Parse("02.01.2006", input)
	if err != nil {
		return "", err
	}
	return date.Format("20060102"), nil
}
