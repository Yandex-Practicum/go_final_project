package base

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

var DateFormat = "20060102"

func CreateDB(envDBFILE string) {
	var appPath string
	var err error
	if envDBFILE != "" {
		appPath = envDBFILE
	} else {
		appPath, err = os.Getwd() //не смогла реализовать через os.Executable()
		if err != nil {
			log.Fatal(err)
		}
	}

	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)

	fmt.Println(err)

	var install bool
	if os.IsNotExist(err) {
		install = true
		fmt.Println("db не найдена")
	} else if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if install {

		// Создание базы данных
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date INTEGER,
			title TEXT NOT NULL DEFAULT "",
			comment TEXT,
			repeat VARCHAR(128) NOT NULL DEFAULT "");`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("База данных успешно создана!")
	}
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		err := errors.New("не указано правило повторения")
		return "", err
	}

	dateForm, err := time.Parse(DateFormat, date)
	if err != nil {
		err := errors.New("указан неверный формат времени")
		return "", err
	}
	rules := strings.Split(repeat, " ")
	var quantity int
	if len(rules) > 1 {
		quantity, err = strconv.Atoi(rules[1])
		if err != nil {
			err = errors.New("неподдерживаемый формат правила повторения")
			return "", err
		}
	}

	switch rules[0] {
	case "d":
		if quantity == 0 {
			err = errors.New("не указан интервал в днях")
			return "", err
		}
		if quantity > 400 {
			err = errors.New("превышен максимально допустимый интервал")
			return "", err
		}
	loopD:
		for {
			dateForm = dateForm.AddDate(0, 0, quantity)
			if dateForm.After(now) {
				break loopD
			}
		}
		return dateForm.Format(DateFormat), nil

	case "y":
	loopY:
		for {
			dateForm = dateForm.AddDate(1, 0, 0)
			if dateForm.After(now) {
				break loopY
			}
		}
		return dateForm.Format(DateFormat), nil

	//case "w":
	//	dateForm = dateForm.AddDate(0, 0, 0)
	//case "m":
	//	dateForm = dateForm.AddDate(0, 0, 0)

	default:
		err = errors.New("недопустимый символ")
		return "", err
	}
}
