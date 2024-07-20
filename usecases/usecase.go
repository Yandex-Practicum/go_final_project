package usecases

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlexJudin/go_final_project/database"
)

var _ Task = (*TaskUsecase)(nil)

type TaskUsecase struct {
	DB database.Task
}

func NewTaskUsecase(db database.Task) *TaskUsecase {
	return &TaskUsecase{DB: db}
}

func (t *TaskUsecase) NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	dateTask, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("error parsing date: %+v", err)
	}

	codeRepeat := strings.Split(repeat, " ")
	if len(codeRepeat) != 2 && !strings.EqualFold(codeRepeat[0], "y") {
		return "", fmt.Errorf("string conversion error repeat")
	}

	switch strings.ToLower(codeRepeat[0]) {
	case "y":
		dateTask = addDateTask(now, dateTask, 1, 0, 0)
	case "d":
		digit, err := convetToIntAncCheck(codeRepeat[1])
		if err != nil {
			return "", err
		}
		dateTask = addDateTask(now, dateTask, 0, 0, digit)
	default:
		return "", fmt.Errorf("unknown key in repeat")
	}

	return dateTask.Format("20060102"), nil
}

func convetToIntAncCheck(num string) (int, error) {
	digit, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}

	if digit >= 400 || digit < 0 {
		return 0, fmt.Errorf("the number of days to reschedule the task is greater than 400 or a negative number is specified")
	}

	return digit, nil
}

func addDateTask(now time.Time, dateTask time.Time, y int, m int, d int) time.Time {
	dateTask = dateTask.AddDate(y, m, d)

	for dateTask.Before(now) {
		dateTask = dateTask.AddDate(y, m, d)
	}

	return dateTask
}
