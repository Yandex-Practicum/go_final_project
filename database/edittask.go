package database

import (
	"errors"
	"fmt"

	"go_final_project/tasks"
)

func (d *Database) EditTask(t tasks.Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat =? WHERE id = ?"

	res, err := d.db.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.Id)
	if err != nil {
		return fmt.Errorf("failed to edit task: %w", err)
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
