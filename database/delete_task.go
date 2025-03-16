package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
)

// DeleteTask удаляет задачу из базы данных по идентификатору
func DeleteTask(db *sqlx.DB, id string) error {
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(deleteQuery, id)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return errors.New("Ошибка выполнения запроса")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка получения результата запроса: %v", err)
		return errors.New("Ошибка получения результата запроса")
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
