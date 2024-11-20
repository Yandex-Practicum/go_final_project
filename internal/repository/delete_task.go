package repository

import (
	"database/sql"
	"errors"
)

func (rep *Repository) DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id=:id"
	res, err := rep.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	if num, _ := res.RowsAffected(); num == 0 {
		return errors.New("задача не найдена")
	}
	return nil
}
