package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pwd/services"

	_ "modernc.org/sqlite"
)

type Dbelement struct {
	Db *sql.DB
}

var DB Dbelement

func ConnectDb() {

	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {

		DB, err := sql.Open("sqlite", dbFile)
		if err != nil {
			fmt.Printf("Ошибка при попытке соединения с базой данный: %s\n", err.Error())
			return
		}
		defer DB.Close()

		createDb := `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT "");

	CREATE INDEX tasks_rules ON scheduler (repeat);`

		_, err = DB.Exec(createDb)
		if err != nil {
			fmt.Printf("Ошибка создания таблицы: %s\n", err.Error())
			return
		}

	}
}

// добавляем задачу в базу данных
func AddTask(task services.Task) (int64, error) {
	res, err := DB.Db.Exec("INSERT INTO tasks (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
