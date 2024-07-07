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
func (dbHandl *Storage) GetTaskByID(id string) (Task, error) {
	var task Task

	// Измененный запрос: выбираются только определенные поля
	row := dbHandl.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Println(err)
		return Task{}, err
	}
	return task, nil
}

// PutTask отправляет SQL запрос на обновление задачи Task, возвращает ошибку в случае неудачи.
func (dbHandl *Storage) PutTask(updateTask Task) error {
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
func (dbHandl *Storage) DeleteTask(id string) error {
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
func (dbHandl *Storage) GetTasksList(search ...string) ([]Task, error) {
	var tasks []Task
	var rows *sql.Rows
	var err error

	switch {
	case len(search) == 0:
		// Измененный запрос: выбираются только определенные поля
		rows, err = dbHandl.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY id LIMIT :limit", sql.Named("limit", rowsLimit))
	case len(search) > 0:
		search := search[0]
		_, err = time.Parse(DateFormat, search)
		if err != nil {
			// Измененный запрос: выбираются только определенные поля
			rows, err = dbHandl.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
				sql.Named("search", search),
				sql.Named("limit", rowsLimit))
			if err != nil {
				return []Task{}, err
			}
			break
		}
		// Измененный запрос: выбираются только определенные поля
		rows, err = dbHandl.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit",
			sql.Named("date", search),
			sql.Named("limit", rowsLimit))
		if err != nil {
			return []Task{}, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		task := Task{}
		// Измененный сканирование: считываются только определенные поля
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
