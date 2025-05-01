package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Task представляет структуру задачи в базе данных
type Task struct {
	ID      int64   `db:"id" json:"id,string"`
	Date    string  `db:"date" json:"date"`
	Title   string  `db:"title" json:"title"`
	Comment *string `db:"comment" json:"comment"`
	Repeat  *string `db:"repeat" json:"repeat"`
}

// MarshalJSON кастомизирует сериализацию Task в JSON
func (t Task) MarshalJSON() ([]byte, error) {
	type Alias Task
	return json.Marshal(&struct {
		*Alias
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}{
		Alias:   (*Alias)(&t),
		Comment: stringOrEmpty(t.Comment),
		Repeat:  stringOrEmpty(t.Repeat),
	})
}

// stringOrEmpty возвращает пустую строку, если указатель nil
func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// TaskOption представляет опцию для фильтрации задач
type TaskOption func(*taskQueryOptions)

type taskQueryOptions struct {
	search string
	date   string
}

// WithSearch добавляет фильтр по подстроке в title или comment
func WithSearch(search string) TaskOption {
	return func(opts *taskQueryOptions) {
		opts.search = search
	}
}

// WithDate добавляет фильтр по точной дате
func WithDate(date string) TaskOption {
	return func(opts *taskQueryOptions) {
		opts.date = date
	}
}

// GetTask возвращает задачу по её ID
func GetTask(id string) (*Task, error) {
	// Создаём контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, что база данных инициализирована
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// Выполняем запрос
	var task Task
	query := `SELECT * FROM scheduler WHERE id = ?`
	err := db.QueryRowxContext(ctx, query, id).StructScan(&task)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// UpdateTask обновляет задачу в базе данных
func UpdateTask(task *Task) error {
	// Создаём контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, что база данных инициализирована
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Выполняем запрос на обновление
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.ExecContext(ctx, query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Проверяем количество затронутых строк
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	return nil
}

// Tasks возвращает список задач с учётом лимита и фильтров
func Tasks(limit int, options ...TaskOption) ([]*Task, error) {
	// Создаём контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Применяем опции
	opts := &taskQueryOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Формируем SQL-запрос
	var query string
	var args []interface{}

	if opts.date != "" {
		// Фильтр по точной дате
		query = `SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		args = []interface{}{opts.date, limit}
	} else if opts.search != "" {
		// Фильтр по подстроке в title или comment
		pattern := "%" + opts.search + "%"
		query = `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		args = []interface{}{pattern, pattern, limit}
	} else {
		// Без фильтров
		query = `SELECT * FROM scheduler ORDER BY date LIMIT ?`
		args = []interface{}{limit}
	}

	// Выполняем запрос
	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	// Сканируем результаты
	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.StructScan(&task)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	// Проверяем ошибки в процессе сканирования
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return tasks, nil
}

// DeleteTask удаляет задачу из базы данных по её ID
func DeleteTask(id string) error {
	// Создаём контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, что база данных инициализирована
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Выполняем запрос на удаление
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// Проверяем количество затронутых строк
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// UpdateDate обновляет только дату задачи в базе данных
func UpdateDate(next string, id string) error {
	// Создаём контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, что база данных инициализирована
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Выполняем запрос на обновление
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.ExecContext(ctx, query, next, id)
	if err != nil {
		return fmt.Errorf("failed to update date: %w", err)
	}

	// Проверяем количество затронутых строк
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
