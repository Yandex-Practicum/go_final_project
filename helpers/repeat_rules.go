package helpers

import (
	"errors"
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

	case repeat == "w" || strings.HasPrefix(repeat, "w "):
		return "", errors.New("wrong format for rule w")
		/*
			days := make([]int, 0, 7)
			daysOfWeek := strings.Split(repeat[2:], ",")

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
			} */

	// TODO: add m rule
	case strings.HasPrefix(repeat, "m "):
		return "", errors.New("wrong format for rule m")
		/*
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

				if dayNumber < 1 || dayNumber > 31 {
					return "", errors.New("days must be between 1 and 31")
				}

				days = append(days, dayNumber)
			}

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
			} else {
				months = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
			}

			var compareTime time.Time
			if now.Before(taskDate) {
				compareTime = taskDate
			} else {
				compareTime = now
			}

			var indexDay int
			for i, day := range days {
				if day > compareTime.Day() {
					indexDay = i
					break
				}
			}

			var indexMonth int
			for i, month := range months {
				if month >= int(compareTime.Month()) {
					indexMonth = i
					break
				}
			}

			nextDate := compareTime
			if len(monthsData) == 2 {
				var tempMonth int
				for {
					if months[indexMonth] > int(compareTime.Month()) {
						tempMonth = 12 - (int(math.Abs(float64((months[indexMonth] - int(compareTime.Month()))))) % 12)
						indexDay = 0
					} else if months[indexMonth] < int(compareTime.Month()) {
						tempMonth = (months[indexMonth] - int(compareTime.Month())) % 12
						indexDay = 0
					} else {
						tempMonth = 0
					}

					difference := (days[indexDay] - int(nextDate.Weekday()) + 7) % 7
					if difference == 0 {
						difference = 7
					}

					nextDate = nextDate.AddDate(0, tempMonth, difference)

					if nextDate.After(now) {
						return nextDate.Format("20060102"), nil
					}

					indexDay = (indexDay + 1) % len(days)
					indexMonth = (indexMonth + 1) % len(months)
				}
			} else {
				for {
					difference := (days[indexDay] - int(nextDate.Weekday()) + 7) % 7
					if difference == 0 {
						difference = 7
					}

					nextDate = nextDate.AddDate(0, 0, difference)

					if nextDate.After(now) {
						return nextDate.Format("20060102"), nil
					}

					indexDay = (indexDay + 1) % len(days)
				}
			}
		*/
	}

	return "", errors.New("repeat rule not found")
}
