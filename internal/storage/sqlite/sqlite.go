package sqlite

import (
	"cactus3d/go_final_project/internal/models"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {

	log.Printf("Storage %s\n", storagePath)
	_, err := os.Stat(storagePath)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	if install {

		stmt, err := db.Prepare(`
			CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY,
			date CHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128));
			CREATE INDEX IF NOT EXISTS idx_sched_date ON scheduler(date);
			`)

		if err != nil {
			return nil, err
		}

		_, err = stmt.Exec()
		if err != nil {
			return nil, err
		}
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddTask(task *models.Task) (int, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetTaskById(id string) (*models.Task, error) {
	query := "SELECT * FROM scheduler WHERE id = ?"
	row := s.db.QueryRow(query, id)
	task := models.Task{}

	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (s *Storage) GetTasks(offset, limit int) ([]models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ? OFFSET ?"

	tasks := []models.Task{}

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return []models.Task{}, nil
	}

	for rows.Next() {
		task := models.Task{}

		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) GetTasksBySearch(search string, offset, limit int) ([]models.Task, error) {
	query := "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ? OFFSET ?"

	tasks := []models.Task{}

	rows, err := s.db.Query(query, search, search, limit, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return []models.Task{}, nil
	}

	for rows.Next() {
		task := models.Task{}

		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) GetTasksByDate(date string, offset, limit int) ([]models.Task, error) {
	query := "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ? OFFSET ?"

	tasks := []models.Task{}

	rows, err := s.db.Query(query, date, limit, offset)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return []models.Task{}, nil
	}

	for rows.Next() {
		task := models.Task{}

		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(task *models.Task) (int64, error) {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Storage) DeleteTaskById(id string) (int64, error) {
	query := `DELETE FROM scheduler WHERE id=?`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
