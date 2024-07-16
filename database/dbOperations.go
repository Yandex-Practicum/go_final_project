package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/LEbauchoir/go_final_project/config"
	"github.com/LEbauchoir/go_final_project/models"
)

func (d *DbHelper) AddTask(t models.Task) (string, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := d.Db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func (d *DbHelper) DeleteTask(taskID int) error {
	query := "DELETE FROM scheduler WHERE id = ?"

	result, err := d.Db.Exec(query, taskID)
	if err != nil {
		return err
	}

	value, err := result.RowsAffected()
	if value == 0 {
		return fmt.Errorf("ошибка выполнения запроса удаления к БД")
	}

	return err
}

func (d *DbHelper) ReadTaskById(id int) (models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?"
	row := d.Db.QueryRow(query, id)

	var task models.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, fmt.Errorf("задача с id %v не найдена", id)
		}
		log.Printf("ошибка выполнения запроса к БД: %v", err)
		return models.Task{}, fmt.Errorf("ошибка выполнения запроса к БД: %v", err)
	}

	return task, nil
}

func (d *DbHelper) UpdateTask(t models.Task) error {
	query := "UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?"
	_, err := d.Db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса к БД: %v", err)

	}
	return nil
}

func (d *DbHelper) TasksShow() ([]models.Task, error) {
	query := fmt.Sprintf("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT %d", config.LimitReturnRows)
	rows, err := d.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса к БД: %v", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк: %v", err)
	}
	return tasks, nil
}

func (d *DbHelper) GetTask(id int) (models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?"
	row := d.Db.QueryRow(query, id)

	var task models.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, fmt.Errorf("задача с id %v не найдена", id)
		}
		log.Printf("ошибка выполнения запроса к БД: %v", err)
		return models.Task{}, fmt.Errorf("ошибка выполнения запроса к БД: %v", err)
	}

	return task, nil
}

func (d *DbHelper) GetMaxID() (int, error) {
	var maxID int
	maxRow := d.Db.QueryRow(`SELECT MAX(id) FROM scheduler`)
	err := maxRow.Scan(&maxID)
	if err != nil {
		return 0, err
	}
	return maxID, nil
}
