package repository

func (r *Repository) DeleteTask(id int) error {

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if row == 0 || err != nil {
		return err
	}
	return nil
}
