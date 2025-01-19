package db

import (
	"database/sql"
	"fmt"
	"log"
)

func (db *DB) DeleteTask(id int) error {
	_, err := db.db.Exec(
		"DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}
	return nil
}
