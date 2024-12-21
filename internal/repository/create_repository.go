package repository

func (r *Repository) CreateScheduler() error {
	query := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL,
			title VARCHAR(32) NOT NULL,
			comment VARCHAR(128) NOT NULL,
			repeat VARCHAR(128) NOT NULL
		);`
	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	queryIndex := "CREATE INDEX IF NOT EXISTS dates_task ON scheduler (date);"
	if _, err := r.db.Exec(queryIndex); err != nil {
		return err
	}

	return nil
}
