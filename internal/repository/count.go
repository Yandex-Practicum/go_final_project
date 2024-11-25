package repository

func (r *Repository) Count() (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM scheduler`
	row := r.db.QueryRow(query)
	err := row.Scan(&count)
	return count, err
}
