package storage

import (
	"database/sql"
	"errors"
	"go_final_project/internal/task"
	"log"
)

const limitTask = 10

func (s *Storage) AddTask(t *task.Task) (int, error) {
	res, err := s.db.Exec(`insert into scheduler (date,title,comment,repeat) values (:date,:title,:comment,:repeat)`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)

	if err != nil {
		log.Println("Не добавилась задача")
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Не удалось получить id", err)
		return 0, err
	}

	return int(id), nil
}

func (s *Storage) GetList() ([]task.Task, error) {
	rows, err := s.db.Query(`select * from scheduler order by date limit :limit`,
		sql.Named("limit", limitTask),
	)

	if err != nil {
		log.Println("Не удалось получить задачи в запросе", err)
		return nil, err
	}

	defer rows.Close()

	var tasks []task.Task

	for rows.Next() {
		t := task.Task{}

		err = rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Println("Не удалось получить задачи в запросе", err)
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if rows.Err() != nil {
		log.Println("Не удалось получить задачи в запросе", err)
		return nil, rows.Err()
	}

	return tasks, nil
}

func (s *Storage) GetTask(id string) (task.Task, error) {
	var t task.Task

	row := s.db.QueryRow(`select * from scheduler where id = :id`,
		sql.Named("id", id),
	)

	err := row.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		log.Println("Не удалось получить задачу по id:", id, err)

		if errors.Is(err, sql.ErrNoRows) {
			return task.Task{}, errors.New(" ")
		}
		return task.Task{}, err
	}

	return t, nil
}

func (s *Storage) ChangeTask(t task.Task) error {
	_, err := s.db.Exec(`update scheduler set date = :date, title = :title, comment = :comment, repeat = :repeat where id = :id`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.Id),
	)

	if err != nil {
		log.Println("не Удалось обновить задачу:", err)
		return err
	}

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	_, err := s.db.Exec(`delete from scheduler where id = :id`,
		sql.Named("id", id))

	if err != nil {
		log.Println("не удалось удалить задачу", err)
		return errors.New("Задача не найденна")
	}

	return nil
}
