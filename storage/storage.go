package storage

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"final/task"
)

const Limit = 50

type DB struct {
	conn *sql.DB
}

func Createdatabase() (DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return DB{conn: nil}, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		log.Fatal(err)
		return DB{conn: nil}, err
	}
	if install {
		createTableSql := `CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`
		_, err = db.Exec(createTableSql)
		if err != nil {
			log.Fatal(err)
			return DB{conn: nil}, err
		}
		return DB{conn: db}, nil
	}
	return DB{conn: db}, nil
}

func (db *DB) Addtasktodb(task task.Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.conn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.New("Ошибка добавления задачи")
	}
	return id, nil
}

func (db *DB) DeleteQuery(id string) string {
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.conn.Exec(deleteQuery, id)
	if err != nil {
		return "Ошибка выполнения запроса"
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "Ошибка получения результата запроса"
	}
	if rowsAffected == 0 {
		return "Запись не найдена"
	}
	return ""
}

func (db *DB) Update(task task.Task) string {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.conn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return "Ошибка выполнения запроса"
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "Ошибка получения результата запроса"
	}

	if rowsAffected == 0 {
		return "Задача не найдена"
	}
	return ""

}

func (db *DB) Findtask(id string) (task.Task, string) {
	var task task.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.conn.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, "Задача не найдена"
		} else {
			return task, "Ошибка выполнения запроса"
		}

	}
	return task, ""
}

func (db *DB) GetTasks() ([]task.Task, error) {
	rows, err := db.conn.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`, sql.Named("limit", Limit))
	if err != nil {
		return nil, errors.New("Ошибка выполнения запроса: ")
	}
	defer rows.Close()

	tasks := make([]task.Task, 0, 0)

	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, errors.New("Ошибка чтения строки: ")
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("Ошибка обработки результата: ")
	}
	return tasks, nil
}

func (db *DB) Updatetask(date string, id string) string {
	updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := db.conn.Exec(updateQuery, date, id)
	if err != nil {
		return "Ошибка обновления задачи"
	}
	return ""
}
