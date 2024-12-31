package repository

import (
	"database/sql"
	"fmt"
	"go_final_project/config"
	"go_final_project/internal/models"
	"log"
	"strconv"
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

func UpdateTaskToDB(task *models.Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	_, err := config.GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDateToDB(date string, task models.Task) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	_, err := config.GetDB().Exec(query, date, task.Id)
	if err != nil {
		return err
	}
	return nil
}

func GetTaskById(strId string) (models.Task, error) {
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("Ошибка преобразования id: %v", err)
		return models.Task{}, fmt.Errorf("некорректный id: %w", err)
	}

	db := config.GetDB()
	task := models.Task{}

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"

	// Выполняем запрос
	row := db.QueryRow(query, id)
	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Задача с id %d не найдена", id)
			return models.Task{}, fmt.Errorf("задача с id %d не найдена", id)
		}
		log.Printf("Ошибка выполнения запроса: %v", err)
		return models.Task{}, fmt.Errorf("ошибка получения задачи: %w", err)
	}
	return task, nil
}

func DeleteTask(strId string) error {
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Printf("Ошибка преобразования id: %v", err)
		return fmt.Errorf("некорректный id: %w", err)
	}
	db := config.GetDB()
	_, err = db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("ошибка удаления запроса: %w", err)
	}
	return nil
}

func GetTasks(search string) ([]models.Task, error) {
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

	// Обрабатываем результаты и маппим их на модель Task
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Printf("Row scan error: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
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
