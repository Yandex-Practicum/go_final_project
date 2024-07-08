package db

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

func DeleteTask(db *sqlx.DB, id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id =?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
