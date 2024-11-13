package repository

func (r *Repository) CreateScheduler() error {

	query := `CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT "");
	
		CREATE INDEX tasks_rules ON scheduler (repeat);`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil

}
