package repository

import "log"

func (repo *Repository) AddTask(date, title, comment, repeat string) (id int64, er string) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	result, err := repo.db.Exec(query, date, title, comment, repeat)
	if err != nil {
		log.Println(err)
		return -1, "ошибка добавление задачи " + err.Error()
	}
	id, err = result.LastInsertId()
	if err != nil {
		return -1, "ошибка получения индекса последнего добавленной задачи " + err.Error()
	}
	return id, ""
}
