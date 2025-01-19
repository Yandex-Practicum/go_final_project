package db

import (
	"database/sql"
	"final-project-bronner/go/models"
	"fmt"
)

func (db *DB) UpdateTask(t models.Task) error {
	sqlResult, err := db.db.Exec(
		"UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat  WHERE id = :id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.Id),
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}

	affectedCnt, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}
	if affectedCnt == 0 {
		return fmt.Errorf("задача с id: %s не найдена", t.Id)
	}
	return nil
}
