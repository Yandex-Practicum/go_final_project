package api

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

const DateFormat = "20060102"

var (
	ErrDateFormat    = errors.New("формат не соответствует " + DateFormat)
	ErrRepeatEmpty   = errors.New("значение repeat не задано")
	ErrRepeatInvalid = errors.New("значение repeat невалидно")
	ErrRepeatDaysMax = errors.New("значение repeat для правила d: максимальное допустимое значение 400")
)

func FormatDate(date time.Time) string {
	return date.Format(DateFormat)
}

func ParseDate(val string) (time.Time, error) {
	date, err := time.Parse(DateFormat, val)
	if err != nil {
		return time.Now(), ErrDateFormat
	}

	return date, nil
}

// Правила повторения задач:
// d <число> — задача переносится на указанное число дней. Максимально допустимое число равно 400;
// y — задача выполняется ежегодно. Этот параметр не требует дополнительных уточнений.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := ParseDate(dstart)
	if err != nil {
		return "", err
	}

	years, days, err := ParseRepeat(repeat)
	if err != nil {
		return "", err
	}

	for {
		date = date.AddDate(years, 0, days)
		if date.After(now) {
			break
		}
	}

	return FormatDate(date), nil
}

func ParseRepeat(repeat string) (years, days int, err error) {
	if repeat == "" {
		return 0, 0, ErrRepeatEmpty
	}

	if repeat == "y" {
		return 1, 0, nil
	}

	r, err := regexp.Compile(`d (\d{1,3})`)
	if err != nil {
		return 0, 0, ErrRepeatInvalid
	}

	matches := r.FindStringSubmatch(repeat)
	if len(matches) < 2 {
		return 0, 0, ErrRepeatInvalid
	}

	days, err = strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, ErrRepeatInvalid
	}

	if days > 400 {
		return 0, 0, ErrRepeatDaysMax
	}

	return 0, days, nil
}
