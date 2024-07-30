package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/go_final_project/model"
)

var _ Task = (*TaskRepo)(nil)

type TaskRepo struct {
	Db *sqlx.DB
}

const limit = 50

func NewNewRepository(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{Db: db}
}

func (r *TaskRepo) CreateTask(task *model.Task) (int64, error) {
	res, err := r.Db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Debugf("Database.CreateTask: %+v", err)

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Debugf("Database.CreateTask: %+v", err)

		return 0, err
	}

	return id, nil
}

func (r *TaskRepo) GetTasks() (model.TasksResp, error) {
	tasks := make([]model.Task, 0)

	res, err := r.Db.Query(`
	SELECT 
    	id,
    	date,
    	title,
    	comment,
    	repeat
    FROM scheduler 
    WHERE date >= $1
	LIMIT $2`,
		time.Now().Format(model.TimeFormat), limit)
	if err != nil {
		log.Debugf("Database.GetTasks: %+v", err)

		return model.TasksResp{Tasks: tasks}, err
	}

	defer res.Close()

	var task model.Task

	for res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Debugf("Database.GetTasks: %+v", err)

			return model.TasksResp{Tasks: tasks}, err
		}

		tasks = append(tasks, task)
	}

	if err = res.Err(); err != nil {
		return model.TasksResp{Tasks: tasks}, err
	}

	return model.TasksResp{Tasks: tasks}, nil
}

func (r *TaskRepo) GetTasksBySearchString(searchString string) (model.TasksResp, error) {
	tasks := make([]model.Task, 0)

	res, err := r.Db.Query(`
	SELECT 
    	id,
    	date,
    	title,
    	comment,
    	repeat
	FROM scheduler 
	WHERE title LIKE $1 OR comment LIKE $1 
	ORDER BY date
	LIMIT $2`,
		"%"+searchString+"%", limit)
	if err != nil {
		log.Debugf("Database.GetTasksBySearchString: %+v", err)

		return model.TasksResp{Tasks: tasks}, err
	}

	defer res.Close()

	var task model.Task

	for res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Debugf("Database.GetTasksBySearchString: %+v", err)

			return model.TasksResp{Tasks: tasks}, err
		}

		tasks = append(tasks, task)
	}

	if err = res.Err(); err != nil {
		return model.TasksResp{Tasks: tasks}, err
	}

	return model.TasksResp{Tasks: tasks}, nil
}

func (r *TaskRepo) GetTasksByDate(searchDate time.Time) (model.TasksResp, error) {
	tasks := make([]model.Task, 0)

	res, err := r.Db.Query(`
	SELECT 
	    id,
    	date,
    	title,
    	comment,
    	repeat
	FROM scheduler 
	WHERE date = $1
	LIMIT $2`,
		searchDate.Format(model.TimeFormat), limit)
	if err != nil {
		log.Debugf("Database.GetTasksByDate: %+v", err)

		return model.TasksResp{Tasks: tasks}, err
	}

	defer res.Close()

	var task model.Task

	for res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Debugf("Database.GetTasksByDate: %+v", err)

			return model.TasksResp{Tasks: tasks}, err
		}

		tasks = append(tasks, task)
	}

	if err = res.Err(); err != nil {
		return model.TasksResp{Tasks: tasks}, err
	}

	return model.TasksResp{Tasks: tasks}, nil
}

func (r *TaskRepo) GetTaskById(id string) (*model.Task, error) {
	var task model.Task

	res, err := r.Db.Query(`
	SELECT 
	    id,
    	date,
    	title,
    	comment,
    	repeat
	FROM scheduler 
	WHERE id = $1`,
		id)
	if err != nil {
		log.Debugf("Database.GetTaskById: %+v", err)

		return nil, err
	}
	defer res.Close()

	if res.Next() {
		err = res.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Debugf("Database.GetTaskById: %+v", err)

			return nil, err
		}
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	if task.Id == "" {
		err = fmt.Errorf("task id %s not found", id)
		log.Debugf("Database.GetTaskById: %+v", err)

		return nil, err
	}

	return &task, nil
}

func (r *TaskRepo) UpdateTask(task *model.Task) error {
	_, err := r.Db.Exec("UPDATE scheduler SET date = $2, title = $3, comment = $4, repeat = $5 WHERE id = $1", task.Id, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Debugf("Database.UpdateTask: %+v", err)

		return err
	}

	return nil
}

func (r *TaskRepo) MakeTaskDone(id string, date string) error {
	res, err := r.Db.Exec("UPDATE scheduler SET date = $2 WHERE id = $1", id, date)
	if err != nil {
		log.Debugf("Database.MakeTaskDone: %+v", err)

		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Debugf("Database.MakeTaskDone: %+v", err)

		return err
	}

	if count == 0 {
		err = fmt.Errorf("task id %s not found", id)
		log.Debugf("Database.MakeTaskDone: %+v", err)

		return err
	}

	return nil
}

func (r *TaskRepo) DeleteTask(id string) error {
	res, err := r.Db.Exec("DELETE FROM scheduler WHERE id = $1", id)
	if err != nil {
		log.Debugf("Database.DeleteTask: %+v", err)

		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Debugf("Database.DeleteTask: %+v", err)

		return err
	}

	if count == 0 {
		err = fmt.Errorf("task id %s not found", id)
		log.Debugf("Database.DeleteTask: %+v", err)

		return err
	}

	return nil
}
