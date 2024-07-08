package db

import (
	"github.com/jmoiron/sqlx"
)

func PutTask(s *sqlx.DB, updTask Task) error {
	_, err := s.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		updTask.Date, updTask.Title, updTask.Comment, updTask.Repeat, updTask.ID)
	if err != nil {
		return err
	}
	return nil
}
