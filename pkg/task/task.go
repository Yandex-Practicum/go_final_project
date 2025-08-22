package task

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {

	dateYear, dateMonth, dateDay := date.Date()
	nowYear, nowMonth, nowDay := now.Date()
	dateTruncated := time.Date(dateYear, dateMonth, dateDay, 0, 0, 0, 0, time.UTC)
	nowTruncated := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.UTC)
	return dateTruncated.After(nowTruncated)
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("repeat rule is empty")
	}

	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart format: %v", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}
	rule := parts[0]

	switch rule {
	case "d":

		if len(parts) != 2 {
			return "", fmt.Errorf("invalid format for 'd' rule: interval required")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid interval for 'd' rule: %v", err)
		}
		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("interval must be between 1 and 400, got %d", interval)
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "y":

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w", "m":

		return "", fmt.Errorf("unsupported repeat rule: %s", rule)

	default:
		return "", fmt.Errorf("unknown repeat rule: %s", rule)
	}

	return date.Format("20060102"), nil
}
