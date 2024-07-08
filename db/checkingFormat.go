package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (task Task) Check() (Task, error) {
	var date time.Time
	var err error

	if len(task.Date) == 0 || strings.ToLower(task.Date) == "today" {
		date = time.Now()
		task.Date = date.Format(DateFormat)
	} else {
		date, err = time.Parse(DateFormat, task.Date)
		if err != nil {
			return Task{}, err
		}
	}

	if !IsValidID(task.ID) {
		return Task{}, fmt.Errorf("некорректный формат ID")
	}

	dateTrunc := date.Truncate(time.Hour * 24)
	nowTrunc := time.Now().Truncate(time.Hour * 24)

	if dateTrunc.Before(nowTrunc) {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return Task{}, err
		}
	}
	return task, nil
}

func IsValidID(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}
