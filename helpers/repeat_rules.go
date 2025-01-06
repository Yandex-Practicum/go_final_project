package helpers

import (
	"errors"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	now = now.Truncate(24 * time.Hour)

	switch {
	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimSpace(repeat[2:]))
		if err != nil {
			return "", err
		}

		if days < 1 || days > 400 {
			return "", errors.New("days must be between 1 and 400")
		}

		nextDate := taskDate.AddDate(0, 0, days)

		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

		return nextDate.Format("20060102"), nil

	case repeat == "y":
		nextDate := taskDate.AddDate(1, 0, 0)

		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		return nextDate.Format("20060102"), nil

	case strings.HasPrefix(repeat, "w "):
		var daysOfWeek []string
		days := make([]int, 0, 7)

		if repeat == "w" {
			return "", errors.New("wrong format for rule w")
		}

		daysOfWeek = strings.Split(repeat[2:], ",")

		if daysOfWeek[0] == "" {
			return "", errors.New("wrong format for rule w")
		}

		for _, day := range daysOfWeek {
			dayOfWeek, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			if dayOfWeek < 1 || dayOfWeek > 7 {
				return "", errors.New("days of week must be between 1 and 7")
			}

			days = append(days, dayOfWeek)
		}
		sort.Ints(days)

		var compareTime time.Time
		if now.Before(taskDate) {
			compareTime = taskDate
		} else {
			compareTime = now
		}

		var index int
		for i, day := range days {
			if day > int(compareTime.Weekday()) {
				index = i
				break
			}
		}

		nextDate := compareTime
		for {
			difference := (days[index] - int(nextDate.Weekday()) + 7) % 7
			if difference == 0 {
				difference = 7
			}

			nextDate := nextDate.AddDate(0, 0, difference)

			if nextDate.After(now) {
				return nextDate.Format("20060102"), nil
			}

			index = (index + 1) % len(days)
		}
	case strings.HasPrefix(repeat, "m "):
		monthsData := strings.Split(repeat[2:], " ")

		if len(monthsData) > 2 || len(monthsData) < 1 {
			return "", errors.New("wrong format for rule m")
		}

		days := make([]int, 0, 31)
		for _, day := range strings.Split(monthsData[0], ",") {
			dayNumber, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			if !((dayNumber > 0 && dayNumber < 32) || (dayNumber == -1 || dayNumber == -2)) {
				return "", errors.New("days must be between 1 and 31")
			}

			days = append(days, dayNumber)
		}
		sort.Ints(days)

		months := make([]int, 0, 12)
		if len(monthsData) == 2 {
			for _, month := range strings.Split(monthsData[1], ",") {
				monthNumber, err := strconv.Atoi(month)
				if err != nil {
					return "", err
				}

				if monthNumber < 1 || monthNumber > 12 {
					return "", errors.New("months must be between 1 and 12")
				}

				months = append(months, monthNumber)
			}
			sort.Ints(months)
		}

		var compareTime time.Time
		if now.Before(taskDate) {
			compareTime = taskDate
		} else {
			compareTime = now
		}

		minusOneContaines := slices.Contains(days, -1)
		minusTwoContains := slices.Contains(days, -2)

		nextDate := compareTime
		if len(months) == 0 {
			for {
				if nextDate.AddDate(0, 0, 1).Month() != nextDate.Month() {
					lastDate, ok := lastDays(nextDate, minusOneContaines, minusTwoContains)
					if ok && lastDate.After(compareTime) {
						return lastDate.Format("20060102"), nil
					}
				}

				nextDate = nextDate.AddDate(0, 0, 1)

				if slices.Contains(days, nextDate.Day()) {
					return nextDate.Format("20060102"), nil
				}
			}
		} else {
			for {
				if nextDate.AddDate(0, 0, 1).Month() != nextDate.Month() && !slices.Contains(months, int(nextDate.Month())) {
					lastDate, ok := lastDays(nextDate, minusOneContaines, minusTwoContains)
					if ok && lastDate.After(compareTime) {
						return lastDate.Format("20060102"), nil
					}
				}

				if !slices.Contains(months, int(nextDate.Month())) {
					nextDate = nextDate.AddDate(0, 1, 0)
					nextDate = nextDate.AddDate(0, 0, -(nextDate.Day() - 1))
					continue
				}

				nextDate = nextDate.AddDate(0, 0, 1)
				

				if slices.Contains(days, nextDate.Day()) {
					return nextDate.Format("20060102"), nil
				}
			}
		}
	}

	return "", errors.New("repeat rule not found")
}

func lastDays(date time.Time, one, second bool) (time.Time, bool) {
	currentYear, currentMonth, _ := date.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, date.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	beforeLastOfMonth := firstOfMonth.AddDate(0, 1, -2)

	if second {
		return beforeLastOfMonth, true
	}

	if one {
		return lastOfMonth, true
	}

	return time.Time{}, false
}
