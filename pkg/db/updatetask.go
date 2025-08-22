package db

import (
	"fmt"
)

// UpdateDate обновляет дату выполнения задачи в базе данных по её ID
func UpdateDate(newDate string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := db.Exec(query, newDate, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи: %w", err)
	}
	return nil
}
