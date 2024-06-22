package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/qavaleria/go_final_project/models"
)

// NextDate вычисляет следующую дату для задачи в соответствии с правилом повторения
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Правило повторения не указано")
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", fmt.Errorf("Неверный формат даты: %v", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("Неверный формат повторения")
	}

	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("Неверный формат повторения для 'd'")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("Неверное кол-во дней")
		}

		for !date.After(now) {
			date = date.AddDate(0, 0, days)
		}
	case "y":
		if len(parts) != 1 {
			return "", errors.New("Неверный формат повторения для 'y'")
		}

		for !date.After(now) {
			nextYear := date.Year() + 1
			if date.Month() == time.February && date.Day() == 29 {
				for !isLeapYear(nextYear) {
					nextYear++
				}
				date = time.Date(nextYear, date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
			} else {
				date = date.AddDate(1, 0, 0)
			}
		}
	default:
		return "", errors.New("Не поддерживаемый формат повторения")
	}

	return date.Format("20060102"), nil
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	if year%4 == 0 {
		if year%100 == 0 {
			return year%400 == 0
		}
		return true
	}
	return false
}

func AddTask(db *sql.DB, task models.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	var date time.Time
	var err error
	if task.Date == "" {
		date = time.Now()
		task.Date = date.Format("20060102")
	} else {
		date, err = time.Parse("20060102", task.Date)
		if err != nil {
			return 0, errors.New("Дата представлена в неправильном формате")
		}
	}

	now := time.Now()
	if date.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			task.Date = nextDate
		}
	}

	if err := ValidateRepeatRule(task.Repeat); err != nil {
		return 0, err
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
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

// ValidateRepeatRule проверяет формат правила повторения
func ValidateRepeatRule(repeat string) error {
	if repeat == "" {
		return nil
	}

	dPattern := regexp.MustCompile(`^d\s\d+$`)
	yPattern := regexp.MustCompile(`^y$`)

	if !dPattern.MatchString(repeat) && !yPattern.MatchString(repeat) {
		return errors.New("правило повторения указано в неправильном формате")
	}

	return nil
}
