package tasks

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
	date = date.Add(time.Hour * time.Duration(days*24))
	for date.Before(now) {
		date = date.Add(time.Hour * time.Duration(days*24))
	}

	return date.Format("20060101"), nil
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

	date = date.AddDate(1, 0, 0)
	for date.Before(now) {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format("20060101"), nil
}

func validRepeatY(rule string) bool {

	rule = strings.TrimSpace(rule)
	return rule == "y"
}
