package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

// sql-запросы
const (
	SearchLimit         = 20
	PostQuery           = "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);"
	RecentTasksQuery    = "SELECT * FROM scheduler ORDER BY date LIMIT :limit;"
	UpdateTaskQuery     = "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id;"
	SearchByStringQuery = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit;"
	SearchByDateQuery   = "SELECT * FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit;"
	SearchByIDQuery     = "SELECT * FROM scheduler WHERE ID = :id;"
	CreateIndexQuery    = "CREATE INDEX indexdate ON scheduler (date);"
	DeleteTaskQuery     = "DELETE FROM scheduler WHERE id = :id;"
	CreateTableQuery    = `CREATE TABLE scheduler (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date VARCHAR(8) NOT NULL, title VARCHAR(255) NOT NULL,
						   comment TEXT NULL DEFAULT "", repeat VARCHAR(255) NOT NULL);`
)

// Создание таблицы и индекса
func createSchedulerTable() {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(CreateTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(CreateIndexQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// Добавление задачи в БД
func postTask(task Task) (id string, err error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return "", err
	}
	defer db.Close()

	res, err := db.Exec(PostQuery, sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", err
	}
	idInt, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	id = strconv.Itoa(int(idInt))
	return
}

// Получение ближайших задач
func getRecentTasks() (Tasks, error) {
	tasks := []Task{}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return Tasks{tasks}, err
	}
	defer db.Close()
	rows, err := db.Query(RecentTasksQuery, sql.Named("limit", SearchLimit))
	if err != nil {
		return Tasks{tasks}, err
	}
	for rows.Next() {
		var tempTask Task
		err := rows.Scan(&tempTask.Id, &tempTask.Date, &tempTask.Title, &tempTask.Comment, &tempTask.Repeat)
		if err != nil {
			tasks = []Task{}
			return Tasks{tasks}, err
		}
		tasks = append(tasks, tempTask)
	}
	return Tasks{tasks}, nil
}

// Поулучение ближайших задач с поиском
func getTasksBySearch(search string) (Tasks, error) {
	tasks := []Task{}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return Tasks{tasks}, err
	}
	defer db.Close()

	date, err := time.Parse("02.01.2006", search)
	var rows *sql.Rows
	if err != nil {
		search = "%" + search + "%"
		rows, err = db.Query(SearchByStringQuery, sql.Named("search", search), sql.Named("limit", SearchLimit))
		if err != nil {
			return Tasks{tasks}, err
		}
	} else {
		dateStr := date.Format(format)
		rows, err = db.Query(SearchByDateQuery, sql.Named("date", dateStr), sql.Named("limit", SearchLimit))
		if err != nil {
			return Tasks{tasks}, err
		}
	}
	for rows.Next() {
		var tempTask Task
		err := rows.Scan(&tempTask.Id, &tempTask.Date, &tempTask.Title, &tempTask.Comment, &tempTask.Repeat)
		if err != nil {
			tasks = []Task{}
			return Tasks{tasks}, err
		}
		tasks = append(tasks, tempTask)
	}
	return Tasks{tasks}, nil
}

// Получения задачи по ID
func getTaskByID(id string) (Task, error) {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return Task{}, err
	}
	defer db.Close()
	var task Task
	row := db.QueryRow(SearchByIDQuery, sql.Named("id", id))
	err = row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

// Обновление задачи
func updateTask(task Task) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()
	res, err := db.Exec(UpdateTaskQuery, sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.Id))
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return ErrTaskNotFound
	}
	return nil
}

// Удаление задачи
func deleteTask(id string) error {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()
	res, err := db.Exec(DeleteTaskQuery, sql.Named("id", id))
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return ErrTaskNotFound
	}
	return nil
}
