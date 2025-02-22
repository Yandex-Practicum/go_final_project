package tasks

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseMonthDays(days string) ([]int, error) {
	daysParts := strings.Split(days, ",")
	var dayInt []int
	for _, part := range daysParts {
		day, err := strconv.Atoi(part)
		if err != nil || day == 0 || day < -31 || day > 31 {
			return nil, errors.New("invalid day")
		}
		dayInt = append(dayInt, day)
	}
	return dayInt, nil
}

func parseMonth(monthStr string) ([]int, error) {
	monthsParts := strings.Split(monthStr, ",")
	var months []int
	for _, part := range monthsParts {
		month, err := strconv.Atoi(part)
		if err != nil || month < 1 || month > 12 {
			return nil, errors.New("invalid month")
		}
		months = append(months, month)
	}
	return months, nil
}

func isMonthAllowed(currentMonth int, months []int) bool {
	for _, month := range months {
		if currentMonth == month {
			return true
		}
	}
	return false
}

func calculateTargetDate(now, taskDate time.Time, days []int, allowedMonths []int) time.Time {
	year, month, _ := taskDate.Date()
	location := taskDate.Location()

	sort.Ints(days)
	var nearestDate *time.Time

	for {
		if len(allowedMonths) > 0 && !isMonthAllowed(int(month), allowedMonths) {
			month++
			if month > 12 {
				month = 1
				year++
			}
			continue
		}

		for _, day := range days {
			var targetDate time.Time

			if day < 0 {
				firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, location)
				lastDayThisMonth := firstDayNextMonth.AddDate(0, 0, -1)
				targetDate = lastDayThisMonth.AddDate(0, 0, day+1)
			} else {
				targetDate = time.Date(year, month, day, 0, 0, 0, 0, location)
				if targetDate.Month() != time.Month(month) {
					continue
				}
			}

			if targetDate.After(now) {
				if nearestDate == nil || targetDate.Before(*nearestDate) {
					nearestDate = &targetDate
				}
			}
		}

		if nearestDate != nil {
			return *nearestDate
		}

		month++
		if month > 12 {
			month = 1
			year++
		}
	}
}

func MonthRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	if len(rules) < 2 || len(rules) > 3 {
		return "", errors.New("wrong format")
	}

	days, err := parseMonthDays(rules[1])
	if err != nil {
		return "", err
	}

	var months []int
	if len(rules) == 3 {
		months, err = parseMonth(rules[2])
		if err != nil {
			return "", err
		}
	}

	for _, day := range days {
		if day < -2 || day == 0 || day > 31 {
			return "", errors.New("invalid day")
		}
	}

	for {

		if len(months) > 0 && !isMonthAllowed(int(taskDate.Month()), months) {
			taskDate = taskDate.AddDate(0, 1, 0)
			continue
		}

		targetDate := calculateTargetDate(now, taskDate, days, months)
		if targetDate.After(now) {
			return targetDate.Format("20060102"), nil
		}

		taskDate = taskDate.AddDate(0, 1, 0)
	}
}
