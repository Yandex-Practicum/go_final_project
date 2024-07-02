package task

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Memonagi/go_final_project/constants"
)

// AddTask добавляет новую задачу в БД
func AddTask(db *sql.DB, task constants.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CheckRepeat проверяет корректность указанного правила повторения
func CheckRepeat(repeat string) error {
	if repeat == "" {
		return nil
	}
	switch string(repeat[0]) {
	case "y":
		return nil
	case "d":
		s := strings.Split(repeat, " ")
		if len(s) != 2 {
			return errors.New("правило повторения указано в неправильном формате")
		} else {
			days, err := strconv.Atoi(s[1])
			if err != nil || days < 1 || days > 400 {
				return errors.New("указано неверное количество дней")
			}
		}
	case "w":
		s := strings.Split(repeat, " ")
		if len(s) != 2 {
			return errors.New("правило повторения указано в неправильном формате")
		} else {
			weekDays := strings.Split(s[1], ",")
			for _, e := range weekDays {
				wDay, err := strconv.Atoi(e)
				if err != nil || wDay < 1 || wDay > 7 {
					return errors.New("указано неверное количество дней")
				}
			}
		}
	default:
		return errors.New("правило повторения указано в неправильном формате")
	}
	return nil
}

// CheckTitle проверяет наличие заголовка
func CheckTitle(title string) (string, error) {
	if len(title) == 0 {
		return "", errors.New("заголовок задачи не может быть пустым")
	}
	return title, nil
}

// CheckDate проверяет корректность указанной даты
func CheckDate(date string) (string, error) {
	now := time.Now()
	if date == "" || date == "today" {
		return now.Format(constants.DateFormat), nil
	} else {
		outDate, err := time.Parse(constants.DateFormat, date)
		if err != nil {
			return "", errors.New("неправильный формат даты")
		}
		if outDate.Before(now) {
			return now.Format(constants.DateFormat), nil
		} else {
			return outDate.Format(constants.DateFormat), nil
		}
	}
}
