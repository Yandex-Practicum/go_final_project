package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {

	dbFile := "scheduler.db"

	// Инициализация базы данных
	err := Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	fmt.Println("База данных успешно инициализирована!")

	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {

	dbFile := "scheduler.db"

	// Инициализация базы данных
	err := Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	fmt.Println("База данных успешно инициализирована!")

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?", limit)
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

	// Проверяем на ошибки после итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {

	dbFile := "scheduler.db"

	// Инициализация базы данных
	err := Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	fmt.Println("База данных успешно инициализирована!")

	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {

	dbFile := "scheduler.db"

	// Инициализация базы данных
	err := Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	fmt.Println("База данных успешно инициализирована!")

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
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func DeleteTask(id string) error {

	dbFile := "scheduler.db"
	// Инициализация базы данных
	err := Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	fmt.Println("База данных успешно инициализирована!")

	// Проверяем, что id не пустой
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	// Выполняем SQL-запрос на удаление задачи
	result, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err // Возвращаем ошибку, если запрос не удался
	}

	// Проверяем, было ли удалено хотя бы одно значение
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // Ошибка при получении количества затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id: %s", id) // Если задача не найдена
	}

	return nil // Успешное удаление
}

func UpdateDate(next string, id string) error {

	dbFile := "scheduler.db"
	// Инициализация базы данных
	err := Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	fmt.Println("База данных успешно инициализирована!")

	query := "UPDATE tasks SET date = ? WHERE id = ?"
	_, err = db.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %w", err)
	}
	return nil
}
