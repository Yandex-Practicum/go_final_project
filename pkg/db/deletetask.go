package db

import (
	"fmt"
)

// Функция для удаления задачи по ID
func DeleteTask(id string) error {
	result, err := db.Exec("DELETE FROM scheduler WHERE id=?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
