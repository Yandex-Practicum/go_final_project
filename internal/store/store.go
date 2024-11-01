package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"todo-server/internal/donetaskrepeat"
	"todo-server/internal/settings"
	"todo-server/internal/tasks"
)

const (
	limit = 20
)

type Store struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Store {
	return Store{db: db}
}

func (s *Store) CreateTask(t tasks.Task) (string, error) {
	var err error
	err = t.CheckTitle()
	if err != nil {
		return "", err
	}
	parseDate, err := t.CheckData()
	if err != nil {
		return "", err
	}
	t.Date, err = t.CheckRepeat(parseDate)
	if err != nil {
		return "", err
	}
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", fmt.Errorf(`{"error":"Ошибка добавления задачи в базу данных"}`)

	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf(`{"error":"Ошибка получения ID добавленной задачи"}`)
	}
	return fmt.Sprintf("%d", id), nil
}

func (s *Store) GetTask(id string) (tasks.Task, error) {
	var t tasks.Task
	if id == "" {
		return tasks.Task{}, fmt.Errorf(`{"error":"нет индификатора задачи"}`)
	}
	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return tasks.Task{}, fmt.Errorf(`{"error":"Задача не найдена"}`)
	}
	return t, nil
}

func (s *Store) UpdateTask(t tasks.Task) error {
	err := t.CheckID()
	if err != nil {
		return err
	}
	err = t.CheckTitle()
	if err != nil {
		return err
	}
	parseDate, err := t.CheckData()
	if err != nil {
		return err
	}
	t.Date, err = t.CheckRepeat(parseDate)
	if err != nil {
		return err
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err = s.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return fmt.Errorf(`{"error":"Ошибка обновления задачи в базе данных"}`)
	}
	return nil
}

func (s *Store) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf(`{"error":"не указан индификатор задачи"}`)
	}
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return fmt.Errorf(`{"error":"не указан индификатор задачи"}`)
	}
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err = s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf(`{"error":"не получается удалить задачу"}`)
	}
	return nil
}

func (s *Store) SearchTask(search string) ([]tasks.Task, error) {
	var t tasks.Task
	var taskList []tasks.Task
	var rows *sql.Rows
	var err error
	if search == "" {
		rows, err = s.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", limit)
	} else if date, error := time.Parse("02.01.2006", search); error == nil {
		query := "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, date.Format(settings.Template), limit)
	} else {
		search = "%%%" + search + "%%%"
		query := "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
		rows, err = s.db.Query(query, search, search, limit)
	}
	if err != nil {
		return []tasks.Task{}, fmt.Errorf(`{"error":"ошибка запроса"}`)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []tasks.Task{}, fmt.Errorf(`{"error":"ошибка сканирования запроса"}`)
		}
		taskList = append(taskList, t)
	}
	if rows.Err() != nil {
		return []tasks.Task{}, fmt.Errorf(`{"error":"ошибка перебра параметров строки"}`)
	}
	if len(taskList) == 0 {
		taskList = []tasks.Task{}
	}

	return taskList, nil
}

func (s *Store) DoneTask(id string) error {
	var task tasks.Task
	if id == "" {
		return fmt.Errorf(`{"error":"не указан индификатор задачи"}`)
	}
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return fmt.Errorf(`{"error":"не указан индификатор задачи"}`)
	}

	row := s.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return fmt.Errorf(`{"error":"Задача не найдена"}`)
	}
	if task.Repeat == "" {
		_, err := s.db.Exec("DELETE FROM scheduler WHERE id=?", task.ID)
		if err != nil {
			return fmt.Errorf(`{"error":"не получается удалить задачу"}`)
		}
	} else {
		next, err := donetaskrepeat.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf(`{"error":"неверное правило повторения"}`)
		}
		task.Date = next
	}
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = s.db.Exec(query, task.Date, task.ID)
	if err != nil {
		return fmt.Errorf(`{"error":"Ошибка обновления даты выполнения задачи"}`)
	}
	fmt.Println(task)
	return nil
}
