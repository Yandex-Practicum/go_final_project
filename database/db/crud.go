package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// task.go содержит функции CRUD для задач Task

var (
	rowsLimit = 15
)

// AddTask отправляет SQL запрос на добавление переданной задачи Task. Возвращает ID добавленной задачи и/или ошибку.
func (dbHandl *DBHandler) AddTask(task Task) (int64, error) {
	var id int64
	res, err := dbHandl.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, _ = res.LastInsertId()
	}
	return id, err
}

// GetTaskByID возвращает задачу Task с указанным ID, или ошибку.
func (dbHandl *DBHandler) GetTaskByID(id string) (Task, error) {
	var task Task

	row := dbHandl.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Println(err)
		return Task{}, err
	}
	return task, nil

}

// PutTask отправляет SQL запрос на обновление задачи Task, возвращает ошибку в случае неудачи.
func (dbHandl *DBHandler) PutTask(updateTask Task) error {
	res, err := dbHandl.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", updateTask.Date),
		sql.Named("title", updateTask.Title),
		sql.Named("comment", updateTask.Comment),
		sql.Named("repeat", updateTask.Repeat),
		sql.Named("id", updateTask.ID))
	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 1 {
		return fmt.Errorf("ошибка при обновление задачи")
	}
	return nil
}

// DeleteTask отправялет SQL запрос на удаление задачи с указанным ID. Возваращает ошибку в случае неудачи.
func (dbHandl *DBHandler) DeleteTask(id string) error {
	_, err := dbHandl.GetTaskByID(id)
	if err != nil {
		return err
	}

	res, err := dbHandl.db.Exec("DELETE FROM scheduler WHERE id= :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected != 1 {
		return fmt.Errorf("при удаление что-то пошло не так")
	}
	return nil
}

// GetTasksList возвращает послдение добавленные задачи []Task, либо последние добавленные задачи подходящие под поисковой запрос search при его наличие.
// Возвращает ошибку, если что-то пошло не так
func (dbHandl *DBHandler) GetTasksList(search ...string) ([]Task, error) {
	var tasks []Task
	var rows *sql.Rows
	var err error

	switch {
	case len(search) == 0:
		rows, err = dbHandl.db.Query("SELECT * FROM scheduler ORDER BY id LIMIT :limit", sql.Named("limit", rowsLimit))
		if err != nil {
			return []Task{}, err
		}
	case len(search) > 0:
		search := search[0]
		_, err = time.Parse(DateFormat, search)
		if err != nil {
			rows, err = dbHandl.db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
				sql.Named("search", search),
				sql.Named("limit", rowsLimit))
			if err != nil {
				return []Task{}, err
			}
			break
		}
		rows, err = dbHandl.db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit",
			sql.Named("date", search),
			sql.Named("limit", rowsLimit))
		if err != nil {
			return []Task{}, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return []Task{}, err
		}
		tasks = append(tasks, task)

	}
	return tasks, nil
}
