package rules

import (
	"errors"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату на основе текущей даты, исходной даты и правила повтора.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("некорректная дата")
	}

	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	rules := strings.Split(repeat, " ")

	switch rules[0] {
	case "d":
		return handleDayRepeat(now, taskDate, rules)
	case "y":
		return handleYearRepeat(now, taskDate, rules)
	case "w":
		return handleWeekRepeat(now, taskDate, rules)
	case "m":
		return handleMonthRepeat(now, taskDate, rules)
	default:
		return "", errors.New("некорректное правило повторения")
	}
}
