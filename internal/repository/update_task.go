package repository

func (r *Repository) UpdateTask(date, title, comment, repeat string, id int) (int64, error) {

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?`
	res, err := r.db.Exec(query, date, title, comment, repeat, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
