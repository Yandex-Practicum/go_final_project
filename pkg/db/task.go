package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в базу данных
func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) 
	          VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %v", err)
	}

	return id, nil
}

// Tasks возвращает список задач с ограничением по количеству
func Tasks(limit int) ([]*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          ORDER BY date ASC 
	          LIMIT ?`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		var idInt int64

		err := rows.Scan(&idInt, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		task.ID = strconv.FormatInt(idInt, 10)
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации: %v", err)
	}

	return tasks, nil
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("неверный формат ID")
	}

	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE id = ?`

	var task Task
	var idResult int64

	err = DB.QueryRow(query, idInt).Scan(&idResult, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %v", err)
	}

	task.ID = strconv.FormatInt(idResult, 10)

	return &task, nil
}

// UpdateTask обновляет существующую задачу

func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	idInt, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный формат ID")
	}

	query := `UPDATE scheduler 
	          SET date = ?, title = ?, comment = ?, repeat = ? 
	          WHERE id = ?`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, idInt)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный формат ID")
	}

	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, idInt)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
