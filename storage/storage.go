package storage

import (
	"database/sql"
	"time"

	"github.com/Arukasnobnes/go_final_project/models"
	"github.com/Arukasnobnes/go_final_project/other"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) GetTaskByID(id string) (models.Task, error) {
	row := s.DB.QueryRow(`SELECT * FROM scheduler WHERE id = ?`, id)
	var task models.Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return task, err
}

func (s *Storage) InsertTask(task models.Task) (int64, error) {
	res, err := s.DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateTask(task models.Task) error {
	res, err := s.DB.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	res, err := s.DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Storage) GetTasks() ([]models.Task, error) {
	rows, err := s.DB.Query(`SELECT * FROM scheduler ORDER BY date ASC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) NextDate(now time.Time, dateStr string, repeatStr string) (string, error) {
	nextDate, err := other.NextDate(now, dateStr, repeatStr)
	if err != nil {
		return "", err
	}

	return nextDate, nil
}
