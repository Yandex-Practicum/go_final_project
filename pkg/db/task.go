package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// Tasks используется для отправления фронтенду данные о текущих задачах, включай поиск
func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	var rows *sql.Rows
	var err error

	if search != "" {
		if t0, perr := time.Parse("02.01.2006", search); perr == nil {
			dateStr := t0.Format("20060102")
			rows, err = db.Query(
				`SELECT id, date, title, comment, repeat
                 FROM scheduler
                 WHERE date = ?
                 ORDER BY date
                 LIMIT ?`,
				dateStr, limit,
			)
		} else {
			like := "%" + search + "%"
			rows, err = db.Query(
				`SELECT id, date, title, comment, repeat
                 FROM scheduler
                 WHERE title LIKE ? OR comment LIKE ?
                 ORDER BY date
                 LIMIT ?`,
				like, like, limit,
			)
		}
	} else {
		rows, err = db.Query(
			`SELECT id, date, title, comment, repeat
             FROM scheduler
             ORDER BY date
             LIMIT ?`,
			limit,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idInt int64
		var t Task
		if scanErr := rows.Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat); scanErr != nil {
			return nil, scanErr
		}
		t.ID = strconv.FormatInt(idInt, 10)
		tasks = append(tasks, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTask возвращает задачу по id
func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, errors.New("no tasks found")
	}
	query := "select * from scheduler where id = ?"

	var task Task

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	query := "update scheduler set date=?, title=?, comment=?, repeat=? where id=?"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// DeleteTask удаляет задачу по id
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

// UpdateDate обновляет дату
func UpdateDate(newDate, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, newDate, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating date`)
	}
	return nil
}
