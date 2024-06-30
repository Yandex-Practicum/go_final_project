package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/qavaleria/go_final_project/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const FormatDate = "20060102"

// NextDate вычисляет следующую дату для задачи в соответствии с правилом повторения
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Правило повторения не указано")
	}

	date, err := time.Parse(FormatDate, dateStr)
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

	return date.Format(FormatDate), nil
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
		task.Date = date.Format(FormatDate)
	} else {
		date, err = time.Parse(FormatDate, task.Date)
		if err != nil {
			return 0, errors.New("Дата представлена в неправильном формате")
		}
	}

	now := time.Now()
	if date.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(FormatDate)
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

// HandleNextDate обработчик для API запроса /api/nextdate
func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(FormatDate, nowStr)
	if err != nil {
		http.Error(w, "Invalid now format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
