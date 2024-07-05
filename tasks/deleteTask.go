package tasks

import (
	"database/sql"
)

func DeleteTask(id string) error {
	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM scheduler WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
