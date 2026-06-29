package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := Db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task db: Задача не найдена")
	}
	return nil
}

func UpdateDate(id string, next string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := Db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task db: Задача не найдена")
	}
	return nil
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := Db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := Db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task db: некорректный id")
	}

	return nil
}
func GetTasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	var rows *sql.Rows
	var err error

	if search == "" {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		rows, err = Db.Query(query, limit)
	} else {
		if parsedTime, errParse := time.Parse("02.01.2006", search); errParse == nil {
			date := parsedTime.Format("20060102")
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
			rows, err = Db.Query(query, date, limit)
		} else {
			query := `SELECT id, date, title, comment, repeat FROM scheduler 
					  WHERE title LIKE ? OR comment LIKE ? 
					  ORDER BY date ASC LIMIT ?`
			searchPattern := "%" + search + "%"
			rows, err = Db.Query(query, searchPattern, searchPattern, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
