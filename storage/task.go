package storage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"go_final_project/task"
)

const (
	formatOfDate = "20060102"
	taskLimit    = 10
)

func (s *Storage) Add(t *task.Task) (int, error) {
	ins, err := s.db.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		log.Println("can't add task:", err)
		return 0, err
	}

	id, err := ins.LastInsertId()
	if err != nil {
		log.Println("can't get task id:", err)
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetList() ([]task.Task, error) {
	rows, err := s.db.Query(
		`SELECT id, date, title, comment, repeat  
		 FROM scheduler
		 ORDER BY date ASC
		 LIMIT :limit`,
		sql.Named("limit", taskLimit),
	)
	if err != nil {
		log.Println("can't get tasks by GetList:", err)
		return nil, err
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		t := task.Task{}

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Println("can't get tasks by GetList:", err)
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		log.Println("can't get tasks by GetList:", rows.Err())
		return nil, rows.Err()
	}

	return tasks, nil
}

func (s *Storage) SearchTasks(search string) ([]task.Task, error) {

	var rows *sql.Rows
	var err error

	date, err := time.Parse("02.01.2006", search)
	if err != nil {
		log.Println("Search by text")

		rows, err = s.db.Query(
			`SELECT id, date, title, comment, repeat 
			FROM scheduler
			WHERE title LIKE :search OR comment LIKE :search
			ORDER BY date LIMIT :limit`,
			sql.Named("search", "%"+search+"%"),
			sql.Named("limit", taskLimit),
		)
	} else {
		log.Println("Search by date")

		target := date.Format(formatOfDate)

		rows, err = s.db.Query(
			`SELECT * 
			   FROM scheduler
			   WHERE date LIKE :target
			   ORDER BY date LIMIT :limit`,
			sql.Named("target", "%"+target+"%"),
			sql.Named("limit", taskLimit),
		)
	}

	if err != nil {
		log.Println("can't find tasks by SearchTasks:", err)
		return nil, err
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		t := task.Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Println("can't find tasks by SearchTasks.", err)
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		log.Println("can't find tasks by SearchTasks:", err)
		return nil, err
	}

	log.Printf("Found %d tasks", len(tasks))

	return tasks, nil
}

func (s *Storage) GetTask(id string) (task.Task, error) {

	row := s.db.QueryRow(
		`SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE id = :id`,
		sql.Named("id", id),
	)

	var t task.Task

	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		log.Println("can't get task by id:", id, err)

		if errors.Is(err, sql.ErrNoRows) {
			return task.Task{}, errors.New(" ")
		}
		return task.Task{}, err
	}

	return t, nil
}

func (s *Storage) Update(t task.Task) error {

	res, err := s.db.Exec(
		`UPDATE scheduler 
		SET date=:date, title= :title, comment= :comment, repeat= :repeat
		WHERE id= :id`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID),
	)

	if err != nil {
		log.Println("can't update task:", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	log.Printf("updated %d tasks", rowsAffected)
	log.Println("update successful")

	return nil
}

func (s *Storage) DeleteTask(id string) error {

	res, err := s.db.Exec(
		`DELETE
		FROM scheduler
		WHERE id = :id`,
		sql.Named("id", id),
	)

	if err != nil {
		log.Println("can't delete task:", err)
		return errors.New("task not found")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	log.Printf("deleted %d tasks", rowsAffected)
	log.Println("delete task successful")

	return nil
}
