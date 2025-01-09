package helpers

import (
	"errors"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate calculates the next occurrence of a date based on a given repetition rule.
// The function takes the current date, a target date, and a repeat rule as input, and returns
// the next date as a string in "YYYYMMDD" format. It also returns an error if any of the inputs
// are invalid or if the repeat rule is not recognized.
//
// Parameters:
// - now: the current date as a time.Time object.
// - date: the target date as a string in "YYYYMMDD" format.
// - repeat: a string representing the repeat rule. It can be one of the following:
//   - "d <days>": for daily repetition, where <days> is an integer between 1 and 400.
//   - "y": for yearly repetition.
//   - "w <days>": for weekly repetition, where <days> is a comma-separated list of integers
//     representing days of the week (1 for Monday to 7 for Sunday).
//   - "m <days> <months>": for monthly repetition, where <days> is a comma-separated list of days
//     (1 to 31, -1 for the last day, -2 for the second last day of the month), and <months> is an
//     optional comma-separated list of months (1 for January to 12 for December).
//
// Returns:
// - A string with the next occurrence date in "YYYYMMDD" format.
// - An error if the input parameters are invalid or if the repeat rule is not recognized.
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
		if repeat == "w" {
			return "", errors.New("wrong format for rule w")
		}

		daysOfWeek := strings.Split(repeat[2:], ",")
		if daysOfWeek[0] == "" {
			return "", errors.New("wrong format for rule w")
		}

		days := make([]int, 0, len(daysOfWeek))
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

		compareTime := now
		if now.Before(taskDate) {
			compareTime = taskDate
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

		compareTime := now
		if now.Before(taskDate) {
			compareTime = taskDate
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

// lastDays returns either the last or second to last day of the month for a given date.
// If 'second' is true, it returns the second to last day of the month.
// If 'one' is true, it returns the last day of the month.
// If neither 'one' nor 'second' is true, it returns an empty time.Time and false.
//
// Parameters:
// - date: the date from which to calculate the last or second last day of the month.
// - one: a boolean indicating whether to return the last day of the month.
// - second: a boolean indicating whether to return the second to last day of the month.
//
// Returns:
// - A time.Time object representing the calculated day of the month.
// - A boolean indicating if a valid day was found and returned.
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
