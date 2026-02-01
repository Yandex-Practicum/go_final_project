package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/JKasus/go_final_project/pkg/constants"
	"github.com/JKasus/go_final_project/pkg/entities"
)

var daysInterval int

//var weeksInterval int

func afterNow(date time.Time, now time.Time) bool {
	return date.After(now)
}

func checkSymbol(symbol string) bool {
	switch symbol {
	case constants.DaySign, constants.YearSign:
		return true
	default:
		return false
	}
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		err := fmt.Errorf("value %s can not be empty", repeat)
		return "", err
	}

	partsRepeat := strings.Split(repeat, " ")
	if len(partsRepeat) > 2 {
		err := errors.New("repeat parameter can not be consists of more than 2 elements")
		return "", err
	}

	if !checkSymbol(partsRepeat[0]) {
		err := fmt.Errorf("value %s can not be used is the rule of repeating", partsRepeat[0])
		return "", err
	}

	if partsRepeat[0] == constants.DaySign {
		if len(partsRepeat) != 2 {
			err := fmt.Errorf("repeat parameter with first value %s can not consists of less than 2 elements", partsRepeat[0])
			return "", err
		}
	} else if partsRepeat[0] == constants.YearSign {
		if len(partsRepeat) != 1 {
			err := fmt.Errorf("repeat parameter with first value %s can not consists of more than 1 element", partsRepeat[0])
			return "", err
		}
	}

	if partsRepeat[0] == constants.DaySign {
		days, err := strconv.Atoi(partsRepeat[1])
		if err != nil {
			err = fmt.Errorf("value %s can not convert to int", partsRepeat[1])
			return "", err
		}
		daysInterval = days
		if daysInterval > constants.DaysMaxValue || daysInterval < constants.DaysMinValue {
			err := fmt.Errorf("%d: ivalid value for days", days)
			return "", err
		}
	}
	//} else if partsRepeat[0] == constants.WeekSign {
	//	weeks, err := strconv.Atoi(partsRepeat[1])
	//	if err != nil {
	//		err = fmt.Errorf("value %s can not convert to int", partsRepeat[1])
	//		return "", err
	//	}
	//	WeeksInterval = weeks
	//	if WeeksInterval > constants.WeeksMaxValue || WeeksInterval < constants.WeeksMinValue {
	//		err := fmt.Errorf("%d: ivalid value for weeks", weeks)
	//		return "", err
	//	}
	//}

	startDate, err := time.Parse(constants.DateFormat, dstart)
	if err != nil {
		err = errors.New("failed to parse start date: " + err.Error())
		return "", err
	}

	switch partsRepeat[0] {
	case constants.DaySign:
		for {
			startDate = startDate.AddDate(0, 0, daysInterval)
			if afterNow(startDate, now) {
				break
			}
		}
	case constants.YearSign:
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if afterNow(startDate, now) {
				break
			}
		}
	}

	return startDate.Format(constants.DateFormat), nil
}

func CheckDate(task *entities.Task) error {
	now := time.Now()
	var next string
	if task.Date == "" {
		task.Date = now.Format(constants.DateFormat)
	}
	t, err := time.Parse(constants.DateFormat, task.Date)
	if err != nil {
		err = errors.New("invalid date: " + task.Date)
		return err
	}
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			err = errors.New("invalid repeat value: " + task.Repeat)
			return err
		}
	}
	if afterNow(now, t) {
		if len(task.Repeat) == 0 || now.Format(constants.DateFormat) == t.Format(constants.DateFormat) {
			task.Date = now.Format(constants.DateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}
