package database

import (
	"errors"
	"fmt"
)

func (d *Database) DeleteTask(id int) error {
	query := "DELETE FROM scheduler WHERE id = ?"

	res, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	var rowsAffected int64
	if rowsAffected, err = res.RowsAffected(); err != nil {
		return fmt.Errorf("failed to determine number of affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}
