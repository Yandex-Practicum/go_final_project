package repository

import (
	"database/sql"
	"errors"
	"final_project/internal/common"
	nextdate "final_project/internal/handlers/next_date"

	"time"
)

func (rep *Repository) DoneTask(id string) error {
	queryOne := "SELECT date, repeat FROM scheduler WHERE id=:id"
	row := rep.db.QueryRow(queryOne, sql.Named("id", id))
	task := common.Task{}
	if err := row.Scan(&task.Date, &task.Repeat); err != nil {
		return err
	}
	if task.Date == "" {
		return errors.New("не указана дата задачи")
	}
	switch task.Repeat {
	case "":
		queryDelete := "DELETE FROM scheduler WHERE id=:id"
		_, err := rep.db.Exec(queryDelete, sql.Named("id", id))
		if err != nil {
			return err
		}
	default:
		date, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
		queryUpdate := "UPDATE scheduler SET date=:date WHERE id=:id"
		_, err = rep.db.Exec(queryUpdate, sql.Named("date", date), sql.Named("id", id))
		if err != nil {
			return err
		}
	}
	return nil
}
