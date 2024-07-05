package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Ikamenev/model"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var db *sql.DB

func InitializationDatabase() {
	_, err := os.Stat("scheduler.db")

	if err != nil {
		db, err = createDbFile("scheduler.db")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = sql.Open("sqlite3", "scheduler.db")
	}
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS `scheduler` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `date` VARCHAR(8) NULL, `title` VARCHAR(64) NOT NULL, `comment` VARCHAR(255) NULL, `repeat` VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}

}

func createDbFile(dbFilePath string) (*sql.DB, error) {
	_, err := os.Create(dbFilePath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertTask(task model.Task) (int, error) {
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func ReadTasks() ([]model.Task, error) {
	var tasks []model.Task

	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date")
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

func SearchTasks(search string) ([]model.Task, error) {
	var tasks []model.Task

	search = fmt.Sprintf("%%%s%%", search)
	rows, err := db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date",
		sql.Named("search", search))
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

func SearchTasksByDate(date string) ([]model.Task, error) {
	var tasks []model.Task

	rows, err := db.Query("SELECT * FROM scheduler WHERE date = :date",
		sql.Named("date", date))
	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []model.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

func ReadTask(id string) (model.Task, error) {
	var task model.Task

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func UpdateTask(task model.Task) (model.Task, error) {
	result, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		return model.Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}

	if rowsAffected == 0 {
		return model.Task{}, errors.New("failed to update")
	}

	return task, nil
}

func DeleteTaskDb(id string) error {
	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete")
	}

	return err
}
