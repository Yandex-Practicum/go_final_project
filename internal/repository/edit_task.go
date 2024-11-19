package repository

import (
	"database/sql"
	"errors"
)

func (rep *Repository) EditTask(id, date, title, comment, repeat string) error {
	query := "UPDATE scheduler SET date=:date, title=:title, comment=:comment,repeat=:repeat WHERE id=:id"
	res, err := rep.db.Exec(query, sql.Named("date", date), sql.Named("title", title), sql.Named("comment", comment), sql.Named("id", id), sql.Named("repeat", repeat))
	if err != nil {
		return err
	}
	if k, _ := res.RowsAffected(); k == 0 {
		return errors.New("изменения не внесены")
	}
	return nil
}
