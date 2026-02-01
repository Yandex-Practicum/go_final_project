package database

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/models"
	"log"
)

// GetTaskByID извлекает задачу из базы данных по идентификатору
func GetTaskByID(db *sqlx.DB, id string) (models.Task, error) {
	var task models.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, errors.New("задача не найдена")
		}
		log.Printf("Ошибка выполнения запроса: %v", err)
		return task, errors.New("Ошибка выполнения запроса")
	}
	return task, nil
}
