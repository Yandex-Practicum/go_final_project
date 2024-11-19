package repository

import "log"

func (rep *Repository) CreateSchedulerIndex() {
	query := `CREATE  INDEX IF NOT EXISTS dateIndex ON scheduler(date)`
	_, err := rep.db.Exec(query)
	if err != nil {
		log.Fatal("ошибка создания индекса ", err)
	}
}
