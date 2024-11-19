package repository

import "database/sql"

func (rep *Repository) DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id=:id"
	_, err := rep.db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}
