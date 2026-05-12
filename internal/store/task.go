package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"go_final_project/internal/models"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskStore struct {
	conn *sql.DB
	logs *log.Logger
}

func NewTaskStore(conn *sql.DB, logs *log.Logger) *TaskStore {
	return &TaskStore{
		conn: conn,
		logs: logs,
	}
}

func (s *TaskStore) List(limit int) ([]models.Task, error) {
	query := "SELECT id, title, comment, date, repeat FROM scheduler LIMIT :limit"

	rows, err := s.conn.Query(query,
		sql.Named("limit", limit))

	if err != nil {
		s.logs.Println(err)

		return nil, err
	}

	defer rows.Close()

	tasks := make([]models.Task, 0, limit)

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Comment,
			&task.Date,
			&task.Repeat)

		if err != nil {
			s.logs.Println(err)

			return nil, err
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err == sql.ErrNoRows {
		s.logs.Println(err)

		return nil, nil
	}

	if err != nil {
		s.logs.Println(err)

		return nil, err
	}

	return tasks, nil
}

func (s *TaskStore) Get(id string) (*models.Task, error) {
	var task models.Task

	query := "SELECT id, title, comment, date, repeat FROM scheduler WHERE id = :id"

	err := s.conn.QueryRow(query, sql.Named("id", id)).Scan(
		&task.ID,
		&task.Title,
		&task.Comment,
		&task.Date,
		&task.Repeat)

	if err == sql.ErrNoRows {
		s.logs.Println(err)

		return nil, ErrTaskNotFound
	}

	if err != nil {
		s.logs.Println(err)

		return nil, err
	}

	return &task, nil
}

func (s *TaskStore) Add(task *models.Task) (string, error) {
	query := "INSERT INTO scheduler (title, comment, date, repeat) VALUES (:title, :comment, :date, :repeat)"

	result, err := s.conn.Exec(query,
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("date", task.Date),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		s.logs.Println(err)

		return "", err
	}

	id, err := result.LastInsertId()

	if err != nil {
		s.logs.Println(err)

		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

func (s *TaskStore) Update(task *models.Task) error {
	query := "UPDATE scheduler SET title = :title, comment = :comment, date = :date, repeat = :repeat WHERE id = :id"

	result, err := s.conn.Exec(query,
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("date", task.Date),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))

	if err != nil {
		s.logs.Println(err)

		return err
	}

	count, err := result.RowsAffected()
	if count == 0 {
		s.logs.Println(ErrTaskNotFound)

		return ErrTaskNotFound
	}

	if err != nil {
		s.logs.Println(err)

		return err
	}

	return nil
}

func (s *TaskStore) UpdateDate(task *models.Task) error {
	query := "UPDATE scheduler SET date = :date WHERE id = :id"

	_, err := s.conn.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("id", task.ID))

	if err == sql.ErrNoRows {
		s.logs.Println(err)

		return ErrTaskNotFound
	}

	if err != nil {
		s.logs.Println(err)

		return err
	}

	return nil
}

func (s *TaskStore) Delete(id string) error {
	query := "DELETE FROM scheduler WHERE id = :id"

	_, err := s.conn.Exec(query, sql.Named("id", id))

	if err == sql.ErrNoRows {
		s.logs.Println(err)

		return ErrTaskNotFound
	}

	if err != nil {
		s.logs.Println(err)

		return err
	}

	return nil
}
