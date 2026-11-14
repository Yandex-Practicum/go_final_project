package api

import (
	"time"

	"go_final_project/pkg/nextdate"
)

const dateFormat = "20060102"

func checkDate(taskDate, repeat string) (string, error) {
	now := time.Now()
	nowDate := now.Format(dateFormat)

	if taskDate == "" {
		taskDate = nowDate
	}

	t, err := time.Parse(dateFormat, taskDate)
	if err != nil {
		return "", err
	}

	today, err := time.Parse(dateFormat, nowDate)
	if err != nil {
		return "", err
	}

	var next string

	if repeat != "" {
		next, err = nextdate.NextDate(now, taskDate, repeat)
		if err != nil {
			return "", err
		}
	}

	if t.Before(today) {
		if repeat == "" {
			taskDate = nowDate
		} else {
			taskDate = next
		}
	}

	return taskDate, nil
}
