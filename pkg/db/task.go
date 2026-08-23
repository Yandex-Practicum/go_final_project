package db

import (
	"fmt"
	"time"
)

// Task представляет задачу планировщика.
type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в таблицу scheduler 
// и возвращает ID добавленной записи.
func AddTask(task *Task) (int64, error) {

	query := `
	INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?) `

	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf(`add task: %w`, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(`get inserted task ID: %w`, err)
	}
	return id, nil
}


// Tasks реализует логику получения данных из БД
// принимает максимальное кол-во записей и возвращает их
func Tasks(limit int, search string) ([]*Task, error) {
	
	query := `
	SELECT id, date, title, comment, repeat
	FROM scheduler
	ORDER BY date
	LIMIT ?`

	args := []any{limit}

	// Если search соответсвует формату даты, ищем задачи на конкретный день.
	// В противном случае выполняем поиск по названию и комментарию.
	if search != "" {
		searchDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date
			LIMIT ?`

			args = []any{
				searchDate.Format("20060102"), 
				limit,
			}

		} else {
			query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date
			LIMIT ?`

			search = "%" + search + "%"
			
			args = []any{search, 
				search, 
				limit,
			}
		}
	}
		

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf(`get tasks: %w`, err)
	}

	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}

		err := rows.Scan(
			&task.ID, 
			&task.Date, 
			&task.Title, 
			&task.Comment, 
			&task.Repeat,
		)
		if err != nil {
			return nil, fmt.Errorf(`scan task: %w`, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(`iterate tasks: %w`, err)
	}
	return tasks, nil
}

// GetTask возвращает задачу по ее идентификатору.
func GetTask(id string) (*Task, error) {
	query := `
	SELECT id, date, title, comment, repeat
	FROM scheduler
	WHERE id = ?`


	task := &Task{}

	err := db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		return nil, fmt.Errorf(`get task: %w`, err)
	}
	return task, nil
}


// UpdateTask обновляет все поля указанной задачи.
func UpdateTask(task *Task) error {
	query := `
	UPDATE scheduler
	SET 
		date = ?,
		title = ?,
		comment = ?,
		repeat = ?
	WHERE id = ?`

	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf(`update task: %w`, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect ID for updating task`)
	}
	return nil
}

// UpdateDate изменяет дату следующего выполнения задачи.
func UpdateDate(next string, id string) error {
	query := `
	UPDATE scheduler
	SET date = ?
	WHERE id = ?`


	res, err := db.Exec(
		query,
		next,
		id,
	)
	if err != nil {
		return fmt.Errorf(`update date: %w`, err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect ID for updating date`)
	}
	return nil
}

// DeleteTask удаляет задачу по ее идентификатору.
func DeleteTask(id string) error {
	query := `
	DELETE FROM scheduler
	WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf(`delete task: %w`, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect ID for deleting task`)
	}
	return nil
}

