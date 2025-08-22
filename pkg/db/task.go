package dbase

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}


func AddTask(task *Task) (int64, error) {
	var id int64
	if DB == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	res, err := DB.Exec("INSERT INTO scheduler (date,title,comment,repeat) VALUES ( :p2, :p3, :p4, :p5)",
		sql.Named("p2", task.Date),
		sql.Named("p3", task.Title),
		sql.Named("p4", task.Comment),
		sql.Named("p5", task.Repeat))

	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return id, nil
}


func Tasks(limit int) ([]*Task, error) {
	var res []*Task
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	rows, err := DB.Query("SELECT *from scheduler ORDER BY date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	row := DB.QueryRow("SELECT *from scheduler WHERE id = :p1", sql.Named("p1", id))

	task := &Task{}

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return task, nil
}


func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	row, err := DB.Exec("UPDATE scheduler SET date = :p1, title = :p2, comment = :p3, repeat = :p4 WHERE id = :p5",
		sql.Named("p1", task.Date),
		sql.Named("p2", task.Title),
		sql.Named("p3", task.Comment),
		sql.Named("p4", task.Repeat),
		sql.Named("p5", task.ID))
	if err != nil {
		return err
	}
	count, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// ***************************************************
func DeleteTask(id string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	result, err := DB.Exec("DELETE from scheduler WHERE id = :p1", sql.Named("p1", id))
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // Ошибка при получении количества строк
	}
	if rowsAffected == 0 {
		return fmt.Errorf("запись не найдена") // Специальная ошибка "не найдено"
	}
	return nil
}

// ***************************************************

func UpdateDate(next string, id string) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	row, err := DB.Exec("UPDATE scheduler SET date = :p1 WHERE id = :p2",
		sql.Named("p1", next),
		sql.Named("p2", id))
	if err != nil {
		return err
	}
	count, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}
