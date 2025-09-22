package repository

import (
	"database/sql"
	"go_final_project/structurs"
	"os"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func NewDBwork(db *sql.DB) DB {
	return DB{db: db}
}
func (s DB) AddDB(DBFile string) error {
	_, err := os.Create(DBFile)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL CHECK(date !=''), title TEXT NOT NULL CHECK(title !=''), comment TEXT, repeat TEXT);")
	if err != nil {
		return err
	}
	_, err = s.db.Exec("CREATE INDEX ID_Date ON scheduler (date);")
	if err != nil {
		return err
	}
	return nil
}
func (s DB) CheckTable() error {
	_, err := s.db.Query("SELECT * FROM scheduler;")
	if err != nil {
		return err
	}
	return nil
}
func (s DB) CreateTable() error {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL CHECK(date !=''), title TEXT NOT NULL CHECK(title !=''), comment TEXT, repeat TEXT);")
	if err != nil {
		return err
	}
	_, err = s.db.Exec("CREATE INDEX ID_Date ON scheduler (date);")
	if err != nil {
		return err
	}
	return nil
}
func (s DB) AddTask(t structurs.Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}
func (s DB) GetTasks() ([]structurs.Tasks, error) {
	var errStr []structurs.Tasks
	rows, err := s.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC;")
	if err != nil {
		return errStr, err
	}
	defer rows.Close()

	var res []structurs.Tasks
	for rows.Next() {
		p := structurs.Tasks{}
		err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			return errStr, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return errStr, err
	}

	return res, nil
}
func (s DB) GetTaskId(id string) (structurs.Tasks, error) {
	p := structurs.Tasks{}
	rows := s.db.QueryRow("SELECT * FROM scheduler WHERE id = :id;",
		sql.Named("id", id))
	err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
	return p, err
}
func (s DB) SearchTask(input string) ([]structurs.Tasks, error) {
	var errStr []structurs.Tasks
	rows, err := s.db.Query("SELECT * FROM scheduler WHERE id LIKE CONCAT('%', :input, '%') OR date LIKE CONCAT('%', :input, '%') OR title LIKE CONCAT('%', :input, '%') OR comment LIKE CONCAT('%', :input, '%') ORDER BY date DESC;",
		sql.Named("input", input))
	if err != nil {
		return errStr, err
	}
	defer rows.Close()

	var res []structurs.Tasks
	for rows.Next() {
		p := structurs.Tasks{}

		err := rows.Scan(&p.Id, &p.Date, &p.Title, &p.Comment, &p.Repeat)
		if err != nil {
			return errStr, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return errStr, err
	}

	return res, nil

}
func (s DB) PutTaskId(t structurs.Tasks) error {
	_, err := s.db.Exec("UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id;",
		sql.Named("id", t.Id),
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return err
	}
	return nil
}
func (s DB) DeleteTaskId(id string) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id;",
		sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

func (s DB) NextDate(now string, date string, repeat string) (string, error) {
	format := "20060102"
	dateN, err := time.Parse(format, date)
	if err != nil {
		return "", err
	}
	nowN, err := time.Parse(format, now)
	if err != nil {
		return "", err
	}

	var yearAdd time.Time
	if (dateN.After(nowN) || dateN == nowN) && repeat == "y" && len(repeat) == 1 {
		yearAdd = dateN.AddDate(1, 0, 0)
	} else if (dateN.After(nowN) || dateN == nowN) && strings.Contains(repeat, "d") && len(repeat) > 2 {
		repeatSplit := strings.Split(repeat, " ")
		day := repeatSplit[1]
		i, err := strconv.Atoi(day)
		if err != nil {
			return "err", err
		}
		if i < 401 && i > 0 {
			Add := dateN.AddDate(0, 0, i)
			dayAdd := Add.Format(format)
			return dayAdd, err
		} else {
			return "err", err
		}
	} else if dateN.Before(nowN) && repeat == "y" && len(repeat) == 1 {
		var yearAddbefore string
		for dateN.Before(nowN) {
			Add := dateN.AddDate(1, 0, 0)
			dateN = Add
			yearAddbefore = Add.Format(format)
		}
		return yearAddbefore, err
	} else if dateN.Before(nowN) && strings.Contains(repeat, "d") && len(repeat) > 2 {
		repeatSplit := strings.Split(repeat, " ")
		day := repeatSplit[1]
		i, err := strconv.Atoi(day)
		if err != nil {
			return "err", err
		}
		if i < 401 && i > 0 {
			var dayAddbefore string
			for dateN.Before(nowN) {
				Add := dateN.AddDate(0, 0, i)
				dateN = Add
				dayAddbefore = Add.Format(format)
			}
			return dayAddbefore, err
		} else {
			return "err", err
		}
	}
	year := yearAdd.Format(format)
	return year, err
}
