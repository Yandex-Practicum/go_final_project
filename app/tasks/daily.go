package tasks

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func DailyRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("invalid format")
	}

	days, err := strconv.Atoi(rules[1])
	if err != nil || days <= 0 || days > 400 {
		return "", errors.New("invalid interval")
	}

	if taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	} else {
		for !taskDate.After(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}
	}
	return taskDate.Format("20060102"), nil
}

func parseWeekdays(days string) ([]time.Weekday, error) {
	dayStrings := strings.Split(days, ",")
	var weekdays []time.Weekday

	for _, dayStr := range dayStrings {
		dayInt, err := strconv.Atoi(dayStr)
		if err != nil || dayInt < 1 || dayInt > 7 {
			return nil, errors.New("invalid week day")
		}
		weekday := time.Weekday((dayInt % 7))
		weekdays = append(weekdays, weekday)
	}
	return weekdays, nil
}

func WeekRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	if len(rules) != 2 {
		return "", errors.New("invalid format")
	}
	weekdays, err := parseWeekdays(rules[1])
	if err != nil {
		return "", err
	}
	for {
		for _, weekday := range weekdays {
			if taskDate.Weekday() == weekday {
				if taskDate.After(now) {
					return taskDate.Format("20060102"), nil
				}
			}
		}
		taskDate = taskDate.AddDate(0, 0, 1)
	}
}
