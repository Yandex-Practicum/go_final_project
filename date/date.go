package date

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Memonagi/go_final_project/constants"
)

// NextDate функция для определения следующей даты в соответствии с правилом (решено без учета повторения по месяцам)
func NextDate(now time.Time, dateString string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("неверный формат правила")
	}

	date, err := time.Parse(constants.DateFormat, dateString)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	repeatSlice := strings.Split(repeat, " ")
	rule := repeatSlice[0]

	switch rule {
	case "d":
		if len(repeatSlice) != 2 {
			return "", errors.New("неверный формат правила")
		}
		days, err := strconv.Atoi(repeatSlice[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("указано неверное количество дней")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) && !date.Equal(now) {
				break
			}
		}
		return date.Format(constants.DateFormat), nil
	case "y":
		if len(repeatSlice) != 1 {
			return "", errors.New("неверный формат правила")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) && !date.Equal(now) {
				break
			}
		}
		return date.Format(constants.DateFormat), nil
	case "w":
		var week []int
		if len(repeatSlice) != 2 {
			return "", errors.New("неверный формат правила")
		}
		wSlice := strings.Split(repeatSlice[1], ",")
		for _, e := range wSlice {
			wDay, err := strconv.Atoi(e)
			if err != nil || wDay < 1 || wDay > 7 {
				return "", errors.New("указано неверное количество дней")
			}
			week = append(week, wDay)
		}
		for i, day := range week {
			if day == 7 {
				week[i] = 0
			}
		}
		date = date.AddDate(0, 0, 1)
		for _, ok := constants.WeekMap[int(date.Weekday())]; ok; _, ok = constants.WeekMap[int(date.Weekday())] {
			for _, e := range week {
				if date.Weekday() == time.Weekday(e) && date.After(now) && !date.Equal(now) {
					return date.Format(constants.DateFormat), nil
				}
			}
			date = date.AddDate(0, 0, 1)
		}
		/*
			case "m":
				var dayOfMonth, month []int
				var newDate time.Time
				switch len(repeatSlice) {
				case 2:
					mSliceDay := strings.Split(repeatSlice[1], ",")
					for _, e := range mSliceDay {
						mDay, err := strconv.Atoi(e)
						if err != nil || mDay < -2 || mDay > 31 {
							return "", errors.New("указано неверное количество дней")
						}
						dayOfMonth = append(dayOfMonth, mDay)
					}
					for _, e := range dayOfMonth {
						switch e {
						case -1:

						case -2:

						default:

						}
					}
				case 3:

				default:
					return "", errors.New("неверный формат правила")
				}
		*/
	default:
		return "", errors.New("неверный формат правила")
	}

	return date.Format(constants.DateFormat), nil
}
