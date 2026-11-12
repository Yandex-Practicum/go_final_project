package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return date.After(now)
}

func nextWeekday(date time.Time, weekdays []int) time.Time {
	for {
		date = date.AddDate(0, 0, 1)

		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		for _, day := range weekdays {
			if weekday == day {
				return date
			}
		}
	}
}

func isMonthDay(date time.Time, days []int) bool {
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for _, day := range days {
		switch day {
		case -1:
			if date.Day() == lastDay {
				return true
			}
		case -2:
			if date.Day() == lastDay-1 {
				return true
			}
		default:
			if date.Day() == day {
				return true
			}
		}
	}

	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", fmt.Errorf("не указано правило повторения задачи")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("указана некорректная дата задачи")
	}

	parts := strings.Fields(repeat)

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("укажите интервал повторения в днях от 1 до 400")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("интервал повторения должен быть от 1 до 400 дней")
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("для ежегодного повторения не нужно указывать дополнительные параметры")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("укажите дни недели для повторения задачи")
		}

		dayStrings := strings.Split(parts[1], ",")
		weekdays := make([]int, 0, len(dayStrings))

		for _, dayString := range dayStrings {
			day, err := strconv.Atoi(dayString)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("дни недели должны быть указаны числами от 1 до 7")
			}

			weekdays = append(weekdays, day)
		}

		for {
			date = nextWeekday(date, weekdays)
			if afterNow(date, now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("укажите дни месяца и, при необходимости, месяцы повторения")
		}

		dayStrings := strings.Split(parts[1], ",")
		days := make([]int, 0, len(dayStrings))

		for _, dayString := range dayStrings {
			day, err := strconv.Atoi(dayString)
			if err != nil || (day < 1 || day > 31) && day != -1 && day != -2 {
				return "", fmt.Errorf("дни месяца должны быть от 1 до 31, либо -1 или -2")
			}

			days = append(days, day)
		}

		var months []int

		if len(parts) == 3 {
			monthStrings := strings.Split(parts[2], ",")
			months = make([]int, 0, len(monthStrings))

			for _, monthString := range monthStrings {
				month, err := strconv.Atoi(monthString)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("месяцы должны быть указаны числами от 1 до 12")
				}

				months = append(months, month)
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			if !afterNow(date, now) {
				continue
			}

			if len(months) > 0 {
				monthMatches := false

				for _, month := range months {
					if int(date.Month()) == month {
						monthMatches = true
						break
					}
				}

				if !monthMatches {
					continue
				}
			}

			if isMonthDay(date, days) {
				break
			}
		}

	default:
		return "", fmt.Errorf("указан неизвестный формат повторения задачи")
	}

	return date.Format(dateFormat), nil
}
