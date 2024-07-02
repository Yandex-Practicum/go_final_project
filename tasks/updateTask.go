package tasks

import (
	"database/sql"
)

func PutTask(updTask Task) error {
	db, err := sql.Open("sqlite3", "./db/scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		updTask.Date, updTask.Title, updTask.Comment, updTask.Repeat, updTask.ID)
	if err != nil {
		return err
	}
	return nil
}
