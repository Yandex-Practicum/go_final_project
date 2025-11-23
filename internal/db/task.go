package db

import (
	"database/sql"
	"errors"
	"fmt"
)

// Структура Task представляет задачу в планировщике.
// Поля соответствуют колонкам таблицы scheduler в базе данных.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

const (
	queryInsertTask = `
		INSERT INTO scheduler
		(date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`
	querySelectTask = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`
	querySelectTasks = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		LIMIT ?
	`
	queryUpdateTask = `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`
	queryUpdateDate = `
		UPDATE scheduler
		SET date = ?
		WHERE id = ?
	`
	queryDeleteTask = `
		DELETE FROM scheduler
		WHERE id = ?
	`
)

// AddTask добавляет новую задачу в базу данных.
// Параметры:
// db - соединение с базой данных;
// task - указатель на структуру Task с данными задачи.
// Возвращает:
// ID вставленной записи (int64) и ошибку (если возникла).
func AddTask(db *sql.DB, task *Task) (int64, error) {
	// Проверяем, что указатель на задачу не равен nil
	if task == nil {
		return 0, errors.New("task cannot be nil")
	}

	// Выполняем SQL-запрос на добавление задачи
	res, err := db.Exec(queryInsertTask, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to execute insert query: %w", err)
	}

	// Получаем ID вновь созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return id, err
}

// GetTask получает задачу из базы данных по её ID.
// Параметры:
// db - соединение с базой данных;
// id - идентификатор задачи.
// Возвращает:
// указатель на структуру Task и ошибку (если возникла).
func GetTask(db *sql.DB, id string) (*Task, error) {
	// Проверяем, что ID не пустой
	if id == "" {
		return nil, errors.New("ID must not be empty")
	}

	// Создаём пустой экземпляр Task для задачи
	var task Task

	// Выполняем запрос и сканируем результат в структуру task
	err := db.QueryRow(querySelectTask, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	// Проверяем, не было ли ошибок при итерации по строкам
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to scan task data: %w", err)
	}

	return &task, nil
}

// GetTasks получает список задач из базы данных с ограничением по количеству.
// Параметры:
// db - соединение с базой данных;
// limit - максимальное количество возвращаемых задач.
// Возвращает:
// слайс указателей на структуры Task и ошибку (если возникла).
func GetTasks(db *sql.DB, limit int) ([]*Task, error) {
	// Проверяем, что limit не равен нулю
	if limit == 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	// Создаём пустой слайс для хранения задач
	var tasks []*Task

	// Выполняем запрос с ограничением на количество записей
	rows, err := db.Query(querySelectTasks, limit)
	if err != nil {
		return nil, err
	}
	// Гарантируем закрытие курсора после завершения работы
	defer rows.Close()

	// Проходим по всем строкам результата
	for rows.Next() {
		// Создаём локальную переменную для новой задачи
		var task Task
		// Сканируем данные текущей строки в структуру task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		// Добавляем задачу в слайс
		tasks = append(tasks, &task)
	}

	// Проверяем, не было ли ошибок при итерации по строкам
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil

}

// UpdateTask обновляет данные задачи в базе данных.
// Параметры:
// db - соединение с базой данных;
// task - указатель на структуру Task с обновлёнными данными.
// Возвращает ошибку, если операция не удалась.
func UpdateTask(db *sql.DB, task *Task) error {
	// Выполняем SQL-запрос на обновление задачи
	res, err := db.Exec(queryUpdateTask, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	// Получаем количество затронутых строк (должно быть 1 для успешного обновления)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected count: %w", err)
	}

	// Если ни одна строка не была обновлена - задача не найдена
	if count == 0 {
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	return nil
}

// UpdateDate обновляет дату задачи в базе данных.
// Параметры:
// db - соединение с базой данных;
// next - новая дата задачи;
// id - идентификатор задачи.
// Возвращает ошибку, если операция не удалась.
func UpdateDate(db *sql.DB, next string, id string) error {
	// Валидация входных данных: ID не должен быть пустым
	if id == "" {
		return errors.New("task ID must not be empty")
	}

	// Выполняем SQL-запрос на обновление даты задачи
	res, err := db.Exec(queryUpdateDate, next, id)
	if err != nil {
		return fmt.Errorf("failed to execute date update query: %w", err)
	}

	// Получаем количество затронутых строк (должно быть 1 для успешного обновления)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected count: %w", err)
	}

	// Если ни одна строка не была обновлена - задача не найдена
	if count == 0 {
		return fmt.Errorf("task with ID %s not found", id)
	}

	return nil
}

// DeleteTask удаляет задачу из базы данных по ID.
// Параметры:
// db - соединение с базой данных;
// id - идентификатор удаляемой задачи.
// Возвращает ошибку, если операция не удалась.
func DeleteTask(db *sql.DB, id string) error {
	// Проверяем, что ID не пустой
	if id == "" {
		return errors.New("task ID must not be empty")
	}

	// Выполняем SQL-запрос на удаление задачи
	res, err := db.Exec(queryDeleteTask, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	// Получаем количество удалённых строк (должно быть 1 для успешного удаления)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected after delete: %w", err)
	}

	// Если ни одна строка не была удалена - задача не найдена
	if count == 0 {
		return fmt.Errorf("no task with ID %s exists in the database", id)
	}

	return nil
}
