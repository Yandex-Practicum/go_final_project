package databases

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"github.com/FunnyFoXD/go_final_project/models"
	"github.com/FunnyFoXD/go_final_project/tests"
)

var path = getPath()

func getPath() string {
	pathDB := os.Getenv("TODO_DBFILE")
	if pathDB == "" {
		pathDB = tests.DBFile
	}

	return pathDB
}

func CreateDB() error {
	var install bool

	_, err := os.Stat(path)
	if err != nil {
		install = true
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("can't open database: %s", err.Error())
	}
	defer database.Close()

	if install {
		query := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT
		);
		CREATE INDEX idx_date ON scheduler(date);
		`

		_, err = database.Exec(query)
		if err != nil {
			return fmt.Errorf("can't create table: %s", err.Error())
		}
	}

	return nil
}

func InsertTask(date, title, comment, repeat string) (int, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return 0, fmt.Errorf("can't open database: %s", err.Error())
	}
	defer database.Close()

	result, err := database.Exec(`INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat))
	if err != nil {
		return 0, fmt.Errorf("can't insert task: %s", err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("can't get last insert id: %s", err.Error())
	}

	return int(id), nil
}

func GetTasks() ([]models.TaskFromDB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %s", err.Error())
	}
	defer database.Close()

	rows, err := database.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 10`)
	if err != nil {
		return nil, fmt.Errorf("can't get tasks: %s", err.Error())
	}
	defer rows.Close()

	var tasks = make([]models.TaskFromDB, 0, 10)
	for rows.Next() {
		var task models.TaskFromDB
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("can't scan task: %s", err.Error())
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while reading rows: %s", err.Error())
	}

	return tasks, nil
}
