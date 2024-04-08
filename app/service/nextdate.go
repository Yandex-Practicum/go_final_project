package service

import (
	"errors"
	"main/app/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errors.New("repeat is empty string")
	}

	dayMatched, _ := regexp.MatchString(`d \d{1,3}`, repeat)
	yearMatched, _ := regexp.MatchString(`y`, repeat)
	weekMatched, _ := regexp.MatchString(`w [1-7]+(,[1-7])*`, repeat)
	monthMatched, _ := regexp.MatchString(`m (\b(0?[1-9]|[1-2][0-9]|3[0-1]|-1|-2)\b|-1|-2)+(,\b(0?[1-9]|[1-2][0-9]|3[0-1])\b|,-1|,-2)* *(\b(0?[1-9]|1[0-2])\b)*(,\b(0?[1-9]|1[0-2])\b)*`, repeat)

	if dayMatched {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}

		if days > 400 {
			return "", errors.New("maximum days count must be 400")
		}

		parsedDate, err := time.Parse(model.DatePattern, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(0, 0, days)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}

		return newDate.Format(model.DatePattern), nil
	} else if yearMatched {
		parsedDate, err := time.Parse(model.DatePattern, date)
		if err != nil {
			return "", err
		}

		newDate := parsedDate.AddDate(1, 0, 0)

		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}

		return newDate.Format(model.DatePattern), nil
	} else if weekMatched {
		parsedDate, err := time.Parse(model.DatePattern, date)
		weekday := int(parsedDate.Weekday())
		if err != nil {
			return "", err
		}

		var newDate time.Time
		var weekdays []int

		for _, weekdayString := range strings.Split(strings.TrimPrefix(repeat, "w "), ",") {
			weekdayInt, _ := strconv.Atoi(weekdayString)
			weekdays = append(weekdays, weekdayInt)
		}

		updated := false
		for _, v := range weekdays {
			if weekday < v {
				newDate = parsedDate.AddDate(0, 0, v-weekday)
				updated = true
				break
			}
		}

		if !updated {
			newDate = parsedDate.AddDate(0, 0, 7-weekday+weekdays[0])
		}

		for newDate.Before(now) || newDate == now {
			weekday = int(newDate.Weekday())

			if weekday == weekdays[0] {
				for _, v := range weekdays {
					if weekday < v {
						newDate = newDate.AddDate(0, 0, v-weekday)
						weekday = int(newDate.Weekday())
					}
				}
			} else {
				newDate = newDate.AddDate(0, 0, 7-weekday+weekdays[0])
			}
		}

		return newDate.Format(model.DatePattern), nil
	} else if monthMatched {

	}

	return "", errors.New("repeat wrong format")
}
