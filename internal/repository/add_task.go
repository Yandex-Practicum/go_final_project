package repository

func (r *Repository) AddTask(date, title, comment, repeat string) (int64, error) {

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, date, title, comment, repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
