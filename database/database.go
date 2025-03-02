package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/AlexeyVilkov/go_final_project/date"
	"github.com/AlexeyVilkov/go_final_project/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// объявляем глобальную переменную
var Db *sqlx.DB

// подключение к БД
func ConnectDB() (*sqlx.DB, error) {

	// проверяем, есть ли файл БД
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

	// создаём файл БД
	if install {
		file, err := os.Create("scheduler.db")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	// настройка подключения к БД
	db, err := sqlx.Open("sqlite3", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// если БД отсутствовала, то создаём таблицу и индекс
	if install {
		// тексты sql-запросов
		createTableSQL := "CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT, title TEXT, comment TEXT, repeat TEXT)"
		createIndexSQL := "CREATE INDEX ixDate ON scheduler (date asc)"

		// создание таблицы
		if _, err = db.Exec(createTableSQL); err != nil {
			fmt.Println("ERROR: create table - ", err)
			return nil, err
		}

		// создание индекса
		if _, err = db.Exec(createIndexSQL); err != nil {
			fmt.Println("ERROR: create index - ", err)
			return nil, err
		}

		// вставка записи
		/*if _, err = db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES('20250218', 'Тест title', 'Тест comment', 'Тест repeat')"); err != nil {
			fmt.Println("ERROR: insert into - ", err)
			return nil, err
		}*/
	}

	// присваиваем глобальной переменной созданную БД
	Db = db

	return db, nil
}

func PostTask(tsk model.Task) (int, error) {
	// запись задачи в БД
	res, err := Db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES(:date, :title, :comment, :repeat)",
		sql.Named("date", tsk.Date),
		sql.Named("title", tsk.Title),
		sql.Named("comment", tsk.Comment),
		sql.Named("repeat", tsk.Repeat))
	if err != nil {
		fmt.Println("ERROR: insert into - ", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("ERROR: не удалось получить последний id - ", err)
		return 0, err
	}

	return int(id), err
}

func ListTasks(cnt int) (model.TasksType, error) {
	// инициализация во избежание nil
	retData := model.TasksType{Tasks: make([]model.Task, 0)}

	// чтение списка задач из БД, количество для вывода = cnt
	rows, err := Db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
		sql.Named("limit", cnt))
	if err != nil {
		fmt.Println("ERROR: select - ", err)
		return retData, err
	}

	defer rows.Close()

	for rows.Next() {
		p := model.Task{}

		err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			return retData, err
		}

		retData.Tasks = append(retData.Tasks, p)
	}
	err = rows.Err()
	if err != nil {
		return retData, err
	}

	return retData, nil
}

func GetTaskByID(id string) (model.Task, error) {
	// инициализация во избежание nil
	retData := model.Task{}

	// чтение из БД задачи по id
	row := Db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&retData.Id, &retData.Date, &retData.Title, &retData.Comment, &retData.Repeat)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("ERROR: scan by select - ", err)
		return retData, err
	}

	// проверка, что задача найдена в БД
	if retData.Id == "" {
		fmt.Println("ERROR: задача не найдена")
		return retData, errors.New("задача не найдена")
	}

	return retData, nil
}

func UpdateTask(tsk model.Task) error {
	// редактирование задачи в БД
	res, err := Db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", tsk.Date),
		sql.Named("title", tsk.Title),
		sql.Named("comment", tsk.Comment),
		sql.Named("repeat", tsk.Repeat),
		sql.Named("id", tsk.Id))
	if err != nil {
		fmt.Println("ERROR: update - ", err)
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		fmt.Println("ERROR: не удалось получить количество обновленных записей - ", err)
		return err
	}

	// проверка, что задачу обновили в БД
	if cnt == 0 {
		fmt.Println("ERROR: обновление не прошло, задача не найдена")
		return errors.New("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	// удаление задачи из БД
	res, err := Db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		fmt.Println("ERROR: update - ", err)
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		fmt.Println("ERROR: не удалось получить количество удаленных записей - ", err)
		return err
	}

	// проверка, что задачу удалили в БД
	if cnt == 0 {
		fmt.Println("ERROR: удаление не прошло, задача не найдена")
		return errors.New("задача не найдена")
	}

	return nil

}

func DoneTask(id string) error {
	// получаем задачу
	task, err := GetTaskByID(id)
	if err != nil {
		fmt.Println("ERROR: GetTaskByID - ", err)
		return err
	}

	// если задача без повторения, то удаляем ее
	if task.Repeat == "" {
		err = DeleteTask(id)
		if err != nil {
			fmt.Println("ERROR: DeleteTask - ", err)
			return err
		}
		return nil
	}

	// если задача с повторением, то переписываем дату след. выполнения
	now := time.Now()

	task.Date, err = date.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		fmt.Println("ERROR: NextDate - ", err)
		return err
	}

	err = UpdateTask(task)
	if err != nil {
		fmt.Println("ERROR: UpdateTask - ", err)
		return err
	}

	return nil
}
