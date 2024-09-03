package steps

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DateForFormat = "20060102"

func NextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(DateForFormat, r.FormValue("now"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nextDate, err := NextDateTask(now, date, repeat)

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte(nextDate))
}

func NextDateTask(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" || repeat == " " {
		err := fmt.Errorf("правило содержит пустую строку: \"%s\"", repeat)
		return "", err
	}

	parseDate, err := time.Parse(DateForFormat, date)
	if err != nil {
		return fmt.Sprintf("неверный формат даты: \"%s\"", date), err
	}

	split := strings.Split(repeat, " ")

	switch split[0] {
	case "d":
		if len(split) > 1 {
			add, err := strconv.Atoi(split[1])
			if err != nil {
				return fmt.Sprintf("\"%s\" не может быть преобразовано в корректный интервал дней", repeat), err
			}

			if add > 400 {
				err := fmt.Errorf("превышен максимально допустимый интервал дней: %d", add)
				return "", err
			}

			if add < 1 {
				err := fmt.Errorf("недопустимое значение интервала дней: %d", add)
				return "", err
			}

			addDate := parseDate.AddDate(0, 0, add)

			for {
				if addDate.After(now) {
					break
				}
				addDate = addDate.AddDate(0, 0, add)
			}

			return addDate.Format(DateForFormat), nil
		}
		err := fmt.Errorf("не указан интервал в днях: \"%s\"", repeat)
		return "", err

	case "y":
		addDate := parseDate.AddDate(1, 0, 0)

		for {
			if addDate.After(now) {
				break
			}
			addDate = addDate.AddDate(1, 0, 0)
		}
		return addDate.Format(DateForFormat), nil

	case "w":
		if len(split) < 2 {
			err := fmt.Errorf("не указан день недели: \"%s\"", split)
			return "", err
		}

		splitWeekDay := strings.Split(split[1], ",")

		if parseDate.Before(now) {
			parseDate = now
		}

		week := make([]int, 0, len(splitWeekDay))
		for i := 0; i < len(splitWeekDay); i++ {

			day, err := strconv.Atoi(splitWeekDay[i])
			if err != nil {
				err := fmt.Errorf("недопустимый символ: \"%s\"", splitWeekDay[i])
				return "", err
			}

			if day > 7 {
				err := fmt.Errorf("недопустимое значение дня недели: \"%s\"", day)
				return "", err
			}

			var wDay int
			if day > int(parseDate.Weekday()) {
				wDay = day - int(parseDate.Weekday())
			} else if parseDate.Weekday().String() == "Sunday" {
				wDay = day
			} else {
				wDay = 7 + day - int(parseDate.Weekday())
			}
			week = append(week, wDay)
		}

		sort.Ints(week)
		return parseDate.AddDate(0, 0, week[0]).Format(DateForFormat), nil

	case "m":
		if parseDate.Before(now) {
			parseDate = now
		}

		splitMonthDay := strings.Split(split[1], ",")
		monthDay := make([]int, 0, len(splitMonthDay))
		for i := 0; i < len(splitMonthDay); i++ {
			day, err := strconv.Atoi(splitMonthDay[i])
			if err != nil {
				err := fmt.Errorf("недопустимый символ: \"%s\"", splitMonthDay[i])
				return "", err
			}

			if day > 31 || day < -2 {
				err := fmt.Errorf("недопустимый день месяца: \"%s\"", day)
				return "", err
			}

			monthDay = append(monthDay, day)
		}

		var sliceDate []time.Time
		var dateNew time.Time
		var new time.Time

		if len(split) < 3 {

			for i := 0; i < len(monthDay); i++ {
				for j := 1; j <= 12; j++ {

					if monthDay[i] < 0 {
						dateNew = time.Date(parseDate.Year(), time.Month(j), 1+monthDay[i], 0, 0, 0, 0, time.UTC)
						if dateNew.Year() == parseDate.Year()-1 {
							dateNew = dateNew.AddDate(1, 0, 0)
						}
					} else {
						dateNewTemp := time.Date(parseDate.Year(), time.Month(j+1), 0, 0, 0, 0, 0, time.UTC)
						if monthDay[i] > dateNewTemp.Day() {
							dateNew = time.Date(parseDate.Year(), time.Month(j+1), monthDay[i], 0, 0, 0, 0, time.UTC)
						} else {
							dateNew = time.Date(parseDate.Year(), time.Month(j), monthDay[i], 0, 0, 0, 0, time.UTC)
						}

					}

					sliceDate = append(sliceDate, dateNew)
				}
			}

			sort.Slice(sliceDate, func(i, j int) bool {
				return sliceDate[i].Before(sliceDate[j])
			})

			i := 0
			for {
				if sliceDate[i].After(parseDate) {
					new = sliceDate[i]
					break

				}

				i++
			}

			return new.Format(DateForFormat), nil
		}
		splitMonth := strings.Split(split[2], ",")
		month := make([]int, 0, len(splitMonth))
		for i := 0; i < len(splitMonth); i++ {
			m, err := strconv.Atoi(splitMonth[i])
			if err != nil {
				err := fmt.Errorf("недопустимый символ: \"%s\"", splitMonth[i])
				return "", err
			}

			if m > 12 {
				err := fmt.Errorf("недопустимый номер месяца года: \"%s\"", m)
				return "", err
			}

			month = append(month, m)
		}

		for i := 0; i < len(monthDay); i++ {
			for j := 0; j < len(month); j++ {

				if monthDay[i] < 0 {
					dateNew = time.Date(parseDate.Year(), time.Month(month[j]+1), 1+monthDay[i], 0, 0, 0, 0, time.UTC)
					if dateNew.Year() == parseDate.Year()-1 {
						dateNew = dateNew.AddDate(1, 0, 0)
					}

				} else {
					dateNew = time.Date(parseDate.Year(), time.Month(month[j]), monthDay[i], 0, 0, 0, 0, time.UTC)
				}

				sliceDate = append(sliceDate, dateNew)
			}

		}

		sort.Slice(sliceDate, func(i, j int) bool {
			return sliceDate[i].Before(sliceDate[j])
		})

		i := 0
		for {
			if sliceDate[i].After(parseDate) {
				new = sliceDate[i]
				break
			}
			i++
		}

		return new.Format(DateForFormat), nil

	default:
		err := fmt.Errorf("недопустимый символ: \"%s\"", split[0])
		return "", err
	}
}
