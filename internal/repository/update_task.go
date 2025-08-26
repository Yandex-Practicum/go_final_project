package repository

import (
	"database/sql"
	"fmt"
)

// Функция для обновления информации о задаче в базе данных
func (r *Repository) UpdateTask(id, date, title, comment, repeat string) error {
	res, err := r.db.Exec("UPDATE scheduler SET Date = :Date, Title = :Title, Comment = :Comment, Repeat = :Repeat WHERE Id = :Id",
		sql.Named("Date", date), //формируем SQL-запрос на обновление строки в таблице
		sql.Named("Title", title),
		sql.Named("Comment", comment),
		sql.Named("Repeat", repeat),
		sql.Named("Id", id))
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	rowsAffected, err := res.RowsAffected() //Проверка количества затронутых строк
	if err != nil {
		return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с id %s не найдена", id)
	}

	return nil
}
