package tasks

import (
	"errors"
	"time"
)

func YearRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	if len(rules) != 1 {
		return "", errors.New("invalid format")
	}
	if taskDate.After(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	} else {
		for !taskDate.After(now) {
			taskDate = taskDate.AddDate(1, 0, 0)
		}
	}
	return taskDate.Format("20060102"), nil
}
