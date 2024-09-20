package tasks

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func (task *Task) Validate() error {

	var errs []string

	if task.Title == "" {
		errs = append(errs, "the Title is blank")
	}

	currentTime := time.Now().Truncate(24 * time.Hour)
	if task.Date == "" {
		task.Date = currentTime.Format("20060102")
	}

	if task.Repeat == "" {
		checkedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			errs = append(errs, "wrong format of the Date")
		} else {
			if checkedDate.Before(currentTime) {
				task.Date = currentTime.Format("20060102")
			}
		}
	} else {
		nextDate, err := NextDate(currentTime, task.Date, task.Repeat)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			task.Date = nextDate
		}
	}

	if len(errs) > 0 {
		errDescription := strings.Join(errs, "; ")
		return fmt.Errorf("failed to validate task.Task: %s", errDescription)
	}

	return nil
}

func NextDate(now time.Time, sourceDate, repeat string) (string, error) {

	date, err := time.Parse("20060102", sourceDate)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %w", err)
	}

	if strings.Contains(repeat, "d") {
		return nextDateWithD(now, date, repeat)
	} else if strings.Contains(repeat, "y") {
		return nextDateWithY(now, date, repeat)
	} else {
		return "", fmt.Errorf("repeat rule is undefined. The rule: \"%s\"", repeat)
	}

}

func errWrongRepeatFormat(ruleType, repeatRule string) error {
	return fmt.Errorf("wrong repeat rule format with \"%s\": rule - \"%s\"", ruleType, repeatRule)
}

func nextDateWithD(now, date time.Time, repeat string) (string, error) {

	if !validRepeatD(repeat) {
		return "", errWrongRepeatFormat("d", repeat)
	}

	days, _ := strconv.Atoi(repeat[2:])
	for date.Before(now) {
		date = date.Add(time.Hour * time.Duration(days*24))
	}

	return date.Format("20060102"), nil
}

func validRepeatD(rule string) bool {

	rule = strings.TrimSpace(rule)
	if valid, _ := regexp.MatchString(`^d\s\d{1,3}$`, rule); !valid {
		return false
	}

	days, err := strconv.Atoi(rule[2:])
	if err != nil || days > 400 {
		return false
	}

	return true
}

func nextDateWithY(now, date time.Time, repeat string) (string, error) {

	if !validRepeatY(repeat) {
		return "", errWrongRepeatFormat("y", repeat)
	}

	for date.Before(now) {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format("20060102"), nil
}

func validRepeatY(rule string) bool {

	rule = strings.TrimSpace(rule)
	return rule == "y"
}
