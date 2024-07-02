package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func idExists(id string) error {
	db := DB
	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&Task{})
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}

func DbExists() bool {
	dbFile := os.Getenv("TODO_DBFILE")
	_, err := os.Stat(dbFile)
	var exists bool
	if err == nil {
		exists = true
	}
	return exists

}

func StartDB() {
	dbFile := os.Getenv("TODO_DBFILE")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetConnMaxLifetime(time.Hour)
	DB = db
}

func InstallDB() {
	dbFile := os.Getenv("TODO_DBFILE")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	installQuery, err := os.ReadFile("internal/db/install.sql")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = db.Exec(string(installQuery))
	if err != nil {
		log.Println(err)
		return
	}
	db.Close()
}

func AddTask(task Task) (int64, error) {
	db := DB
	var id int64
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, _ = res.LastInsertId()
	}
	return id, err
}

func GetTaskList(search ...string) ([]Task, error) {
	db := DB
	var rowsLimit int = 15
	var tasks []Task
	var rows *sql.Rows
	var err error

	switch {
	case len(search) == 0:
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY id LIMIT :limit", sql.Named("limit", rowsLimit))
	case len(search) > 0:
		search := search[0]
		_, err = time.Parse(Format, search)
		if err != nil {
			rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
				sql.Named("search", search),
				sql.Named("limit", rowsLimit))
			break
		}
		rows, err = db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit",
			sql.Named("date", search),
			sql.Named("limit", rowsLimit))
	}
	if err != nil {
		return []Task{}, err
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

func GetTaskByID(id string) (Task, error) {
	var task Task
	db := DB
	err := idExists(id)
	if err != nil {
		return Task{}, err
	}

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Println(err)
		return Task{}, err
	}
	return task, nil

}

func PutTask(updTask Task) error {
	db := DB

	err := idExists(updTask.ID)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", updTask.Date),
		sql.Named("title", updTask.Title),
		sql.Named("comment", updTask.Comment),
		sql.Named("repeat", updTask.Repeat),
		sql.Named("id", updTask.ID))
	if err != nil {
		return err
	}
	return nil
}

func DeleteTask(id string) error {
	db := DB

	err := idExists(id)
	if err != nil {
		return err
	}

	res, err := db.Exec("DELETE FROM scheduler WHERE id= :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected != 1 {
		return fmt.Errorf("при удаление что-то пошло не так")
	}
	return nil
}
