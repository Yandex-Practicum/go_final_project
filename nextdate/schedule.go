package nextdate

import (
	"errors"
	"fmt"
	"go_final_project/constants"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse(constants.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}
	ruleParts := strings.Split(repeat, " ")
	switch ruleParts[0] {
	case "d":
		if len(ruleParts) != 2 {
			return "", errors.New("invalid repeat format for 'd'")
		}
		days, err := strconv.Atoi(ruleParts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid days in repeat rule")
		}
		for nextDate := startDate.AddDate(0, 0, days); ; nextDate = nextDate.AddDate(0, 0, days) {
			if nextDate.After(now) {
				return nextDate.Format(constants.DateFormat), nil
			}
		}
	case "y":
		for nextDate := startDate.AddDate(1, 0, 0); ; nextDate = nextDate.AddDate(1, 0, 0) {
			if nextDate.After(now) {
				return nextDate.Format(constants.DateFormat), nil
			}
		}
	default:
		return "", errors.New("invalid or unsupported repeat rule")
	}
}

func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
