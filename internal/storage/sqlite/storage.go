package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wissio/go_final_project/internal/models"
	"github.com/wissio/go_final_project/internal/services"
)

const limitTask = 10

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) CreateTask(t *models.Task) (int, error) {
	const (
		op    = "storage.sqlite.CreateTask"
		query = `
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES ($1, $2, $3, $4)
	`
	)
	result, err := s.db.Exec(query,
		t.Date,
		t.Title,
		t.Comment,
		t.Repeat,
	)
	if err != nil {
		err = wrapError(op, err)
		log.Printf("Error: %v\n", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		err = wrapError(op, err)
		log.Printf("Error: %v\n", err)
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetTask(id string) (models.Task, error) {
	const (
		op    = "storage.sqlite.GetTask"
		query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	)
	var task models.Task
	row := s.db.QueryRow(query, id)
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		err = wrapError(op, err)
		if errors.Is(err, sql.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}
	return task, nil
}

func (s *Storage) UpdateTask(task models.Task) error {
	const (
		op    = "storage.sqlite.UpdateTask"
		query = `UPDATE scheduler 
		         SET date = ?, title = ?, comment = ?, repeat = ? 
		         WHERE id = ?`
	)

	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		err = wrapError(op, err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = wrapError(op, err)
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil

}

func (s *Storage) DeleteTask(id int64) error {
	const (
		op    = "storage.sqlite.DeleteTask"
		query = `DELETE FROM scheduler WHERE id = ?`
	)
	res, err := s.db.Exec(query, id)
	if err != nil {
		err = wrapError(op, err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		err = wrapError(op, err)
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (s *Storage) DoneTask(id int64) error {
	const (
		op          = "storage.sqlite.DoneTask"
		selectQuery = "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
		updateQuery = "UPDATE scheduler SET date = ? WHERE id = ?"
		deleteQuery = "DELETE FROM scheduler WHERE id = ?"
	)

	var task models.Task

	err := s.db.QueryRow(selectQuery, id).Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		err = wrapError(op, err)
		if err == sql.ErrNoRows {
			return ErrTaskNotFound
		}
		err = wrapError(op, err)
		return err
	}
	if task.Repeat != "" {
		nextDate, err := services.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			err = wrapError(op, err)
			return err
		}
		task.Date = nextDate
		_, err = s.db.Exec(updateQuery, task.Date, task.Id)
		if err != nil {
			err = wrapError(op, err)
			return err
		}
	} else {
		_, err = s.db.Exec(deleteQuery, task.Id)
		if err != nil {
			err = wrapError(op, err)
			return err
		}
	}
	return nil
}

func (s *Storage) GetTasks(date, search string, limit int) ([]models.Task, error) {
	const (
		op           = "storage.sqlite.GetTasks"
		baseQuery    = `SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1`
		dateFilter   = ` AND date = ?`
		searchFilter = ` AND (title LIKE ? OR comment LIKE ?)`
		orderBy      = ` ORDER BY date`
		limitClause  = ` LIMIT ?`
	)

	var (
		query string
		args  []interface{}
	)

	query = baseQuery

	if date != "" {
		query += dateFilter
		args = append(args, date)
	}

	if search != "" {
		query += searchFilter
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	query += orderBy

	if limit > 0 {
		query += limitClause
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, wrapError(op, err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var id int64
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, wrapError(op, err)
		}
		task.Id = fmt.Sprintf("%d", id)
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, wrapError(op, err)
	}

	return tasks, nil
}

func wrapError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
