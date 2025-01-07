package databases

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"time"

	_ "modernc.org/sqlite"

	"github.com/FunnyFoXD/go_final_project/models"
	"github.com/FunnyFoXD/go_final_project/tests"
)

var path = getPath()

// getPath returns a path to the SQLite database file.
//
// If TODO_DBFILE environment variable is set, it will be used, otherwise
// tests.DBFile will be used as a default value.
func getPath() string {
	pathDB := os.Getenv("TODO_DBFILE")
	if pathDB == "" {
		pathDB = tests.DBFile
	}

	return pathDB
}

// CreateDB creates the database file if it doesn't exist and installs the
// scheduler table schema with an index on the date column.
//
// If the database file already exists, this function does nothing and returns
// nil.
//
// The table schema is as follows:
//
// CREATE TABLE scheduler (
// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
// 	date TEXT NOT NULL,
// 	title TEXT NOT NULL,
// 	comment TEXT,
// 	repeat TEXT
// );
// CREATE INDEX idx_date ON scheduler(date);
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

// InsertTask inserts a new task into the scheduler table.
//
// The function takes four string parameters for the date, title, comment, and
// repeat fields of the task. The date must be in the format "YYYYMMDD".
//
// The function returns the integer ID of the newly inserted task and an error.
// If the task is successfully inserted, the returned error is nil.
//
// If the task cannot be inserted, the function returns an error with the
// following format: "can't insert task: <error message>".
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

// GetTasks returns a slice of models.TaskFromDB and an error.
//
// The function takes one string parameter for a search query.
//
// If the search query is a date in the format "DD.MM.YYYY", the function
// returns a slice of tasks for that date. Otherwise, the function returns a
// slice of tasks whose title or comment contains the search query.
//
// The returned slice of tasks is limited to 20 tasks.
//
// If the tasks cannot be retrieved, the function returns an error with the
// following format: "can't get tasks: <error message>".
func GetTasks(search string) ([]models.TaskFromDB, error) {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %s", err.Error())
	}
	defer database.Close()

	dateRegExp := regexp.MustCompile(`^([0-2][0-9]|(3)[0-1])\.(0[1-9]|1[0-2])\.\d{4}$`)
	matched := dateRegExp.MatchString(search)

	var tasks = make([]models.TaskFromDB, 0, 20)
	var rows *sql.Rows

	switch matched {
	case true:
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			return nil, fmt.Errorf("error while parsing date: %v", err)
		}

		dateStr := date.Format("20060102")

		rows, err = database.Query(`SELECT id, date, title, comment, repeat FROM scheduler
			WHERE date = :date ORDER BY date ASC LIMIT 20`,
			sql.Named("date", dateStr))
		if err != nil {
			return nil, fmt.Errorf("can't get tasks: %v", err)
		}
	default:
		rows, err = database.Query(`SELECT id, date, title, comment, repeat FROM scheduler
			WHERE title LIKE :search OR comment LIKE :search ORDER BY date ASC LIMIT 20`,
			sql.Named("search", "%"+search+"%"))
		if err != nil {
			return nil, fmt.Errorf("can't get tasks: %v", err)
		}
	}
	defer rows.Close()

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

// GetTaskByID gets a task from the database by its id.
//
// The function takes a string id parameter and returns a TaskFromDB struct and an error.
// If the task is not found, the error is sql.ErrNoRows.
// If there is a database error, the error is wrapped with a message in the format
// "can't get task: <error message>".
func GetTaskByID(id string) (models.TaskFromDB, error) {
	var task models.TaskFromDB

	database, err := sql.Open("sqlite", path)
	if err != nil {
		return task, fmt.Errorf("can't open database: %v", err)
	}
	defer database.Close()

	err = database.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`,
		sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err == sql.ErrNoRows {
		return task, err
	} else if err != nil {
		return task, fmt.Errorf("can't get task: %v", err)
	}

	return task, nil
}

// UpdateTaskByID updates an existing task in the scheduler table.
//
// The function takes a TaskFromDB struct as a parameter, which contains the
// updated details of the task including id, date, title, comment, and repeat.
//
// The function returns an error if the task cannot be found or updated.
// If the task is not found, the error is sql.ErrNoRows.
// If there is a database error, the error is wrapped with a message in the
// format "can't update task: <error message>".
func UpdateTaskByID(updatedTask models.TaskFromDB) error {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("can't open database: %v", err)
	}
	defer database.Close()

	var existingID string

	err = database.QueryRow(`SELECT id FROM scheduler WHERE id = :id`,
		sql.Named("id", updatedTask.ID)).Scan(&existingID)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return fmt.Errorf("can't get task: %v", err)
	}

	query := `
	UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat
	WHERE id = :id
	`

	_, err = database.Exec(query,
		sql.Named("id", updatedTask.ID),
		sql.Named("date", updatedTask.Date),
		sql.Named("title", updatedTask.Title),
		sql.Named("comment", updatedTask.Comment),
		sql.Named("repeat", updatedTask.Repeat))
	if err != nil {
		return fmt.Errorf("can't update task: %v", err)
	}

	return nil
}

// UpdateTaskDateByID updates the date field of a task with the given id.
//
// The function returns an error if the task cannot be found or updated.
// If the task is not found, the error is sql.ErrNoRows.
// If there is a database error, the error is wrapped with a message in the
// format "can't update task: <error message>".
func UpdateTaskDateByID(id, date string) error {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("can't open database: %v", err)
	}
	defer database.Close()

	_, err = database.Exec(`UPDATE scheduler SET date = :date WHERE id = :id`,
		sql.Named("id", id),
		sql.Named("date", date))
	if err != nil {
		return fmt.Errorf("can't update task: %v", err)
	}

	return nil
}

// DeleteTask deletes a task with the given id from the database.
//
// The function returns an error if the task cannot be found or deleted.
// If the task is not found, the error is sql.ErrNoRows.
// If there is a database error, the error is wrapped with a message in the
// format "can't delete task: <error message>".
func DeleteTask(id string) error {
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("can't open database: %v", err)
	}
	defer database.Close()

	_, err = database.Exec(`DELETE FROM scheduler WHERE id = :id`,
		sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("can't delete task: %v", err)
	}

	return nil
}
