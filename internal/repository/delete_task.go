package repository

import (
	"database/sql"
	"fmt"
)

// Функция DeleteTask для удаления записи о задаче из базы данных по указанному идентификатору
func (r *Repository) DeleteTask(id string) error {
	res, err := r.db.Exec("DELETE FROM scheduler WHERE Id = :Id", sql.Named("Id", id)) //SQL-запрос на удаление строки
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	rowsAffected, err := res.RowsAffected() //Проверка количества затронутых строк
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с id %s не найдена", id)
	}

	return nil
}
