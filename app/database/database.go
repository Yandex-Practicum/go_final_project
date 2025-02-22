package database

import (
	"database/sql"
	"time"

	"go_final/app/models"
)

func InsertIntoDB(task models.Remind) (uint64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func FindInDB(search string, limit int) ([]models.Remind, error) {
	var query string
	var args []interface{}

	if search == "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
		args = append(args, limit)
	} else {
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?"
			args = append(args, parsedDate.Format("20060102"), limit)
		} else {
			likePattern := "%" + search + "%"
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			args = append(args, likePattern, likePattern, limit)
		}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Remind
	for rows.Next() {
		var task models.Remind
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTaskByID(id int) (models.Remind, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := DB.QueryRow(query, id)
	var task models.Remind
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return models.Remind{}, err
	}
	return task, nil
}

func UpdateTask(task models.Remind) error {
	query := `
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteTaskByID(id int) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := DB.Exec(query, id)
	return err
}

func UpdateTaskDate(id uint64, newDate string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	_, err := DB.Exec(query, newDate, id)
	return err
}
