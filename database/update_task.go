package database

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/models"
	"log"
)

// UpdateTask обновляет задачу в базе данных
func UpdateTask(db *sqlx.DB, task models.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return errors.New("ошибка выполнения запроса")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка получения результата запроса: %v", err)
		return errors.New("ошибка получения результата запроса")
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
