package services

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

var repeatFunctions = map[byte]func(now time.Time, start time.Time, repeat string) (string, error){
	'y': repeatAnnually,
	'd': repeatDaily,
	'w': repeatWeekly,
	'm': repeatMonthly,
}

func NextDate(now time.Time, startDate, repeat string) (string, error) {
	parsedDate, err := time.Parse(layout, startDate)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	if repeat == "" {
		return "", errors.New("empty repeat pattern")
	}

	repeatFunction := repeatFunctions[repeat[0]]
	if repeatFunction == nil {
		return "", fmt.Errorf("invalid repeat pattern: %s", repeat)
	}

	return repeatFunction(now, parsedDate, repeat)
}

func repeatAnnually(now time.Time, startDate time.Time, repeat string) (string, error) {
	nextAnnualDate := startDate.AddDate(1, 0, 0)
	for !nextAnnualDate.After(now) {
		nextAnnualDate = nextAnnualDate.AddDate(1, 0, 0)
	}
	return nextAnnualDate.Format(layout), nil
}

func repeatDaily(now time.Time, startDate time.Time, repeat string) (string, error) {
	repeatPatternParts := strings.SplitN(repeat, " ", 2)
	if len(repeatPatternParts) != 2 || repeatPatternParts[0] != "d" {
		return "", fmt.Errorf("invalid format: %s", repeat)
	}
	daysStr := repeatPatternParts[1]
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 400 {
		return "", fmt.Errorf("invalid format: %s", repeat)
	}

	nextDate := startDate.AddDate(0, 0, days)
	for !nextDate.After(now) {
		nextDate = nextDate.AddDate(0, 0, days)
	}

	return nextDate.Format(layout), nil
}

func repeatWeekly(now time.Time, startDate time.Time, repeat string) (string, error) {
	repeat = strings.TrimPrefix(repeat, "w ")
	daysOfWeekStr := strings.Split(repeat, ",")
	daysOfWeek := make([]time.Weekday, len(daysOfWeekStr))
	for i, dayStr := range daysOfWeekStr {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("invalid format: %s", repeat)
		}
		daysOfWeek[i] = time.Weekday(day % 7)
	}
	sort.Slice(daysOfWeek, func(i, j int) bool {
		return daysOfWeek[i] < daysOfWeek[j]
	})
	next := findNextWeekday(startDate, daysOfWeek)
	for !next.After(now) {
		next = findNextWeekday(next.AddDate(0, 0, 1), daysOfWeek)
	}
	return next.Format(layout), nil
}

func repeatMonthly(now time.Time, startDate time.Time, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat[1:])
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 || len(parts) > 2 {
		return "", fmt.Errorf("invalid format: %s", repeat)
	}
	daysPart := strings.Split(parts[0], ",")
	daysMap := make(map[int]bool)
	for _, day := range daysPart {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt < -2 || dayInt == 0 || dayInt > 31 {
			return "", fmt.Errorf("invalid format: %s", day)
		}
		daysMap[dayInt] = true
	}
	monthsMap := make(map[int]bool)
	if len(parts) == 2 {
		for _, m := range strings.Split(parts[1], ",") {
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("invalid format: %s", parts[1])
			}
			monthsMap[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			monthsMap[i] = true
		}
	}
	for next := startDate; ; next = next.AddDate(0, 0, 1) {
		day := next.Day()
		month := int(next.Month())
		if daysMap[day] || daysMap[day-daysInMonth(next.Month(), next.Year())-1] {
			if monthsMap[month] && next.After(now) {
				return next.Format(layout), nil
			}
		}
		if next.Year() > now.Year()+1 {
			break
		}
	}
	return "", errors.New("no next suitable date found")
}

func findNextWeekday(start time.Time, weekdays []time.Weekday) time.Time {
	currentDay := start.Weekday()
	weekdaysSet := make(map[time.Weekday]bool)
	for _, day := range weekdays {
		weekdaysSet[day] = true
	}

	daysUntilNextWeekday := int(weekdays[0] - currentDay)
	if daysUntilNextWeekday <= 0 {
		daysUntilNextWeekday += 7
	}

	for i := 0; i < 7; i++ {
		nextDay := start.AddDate(0, 0, daysUntilNextWeekday)
		if weekdaysSet[nextDay.Weekday()] {
			return nextDay
		}
		daysUntilNextWeekday++
		if daysUntilNextWeekday == 7 {
			daysUntilNextWeekday = 0
		}
	}

	return time.Time{}
}

func daysInMonth(month time.Month, year int) int {
	switch month {
	case time.February:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	case time.April, time.June, time.September, time.November:
		return 30
	default:
		return 31
	}
}
