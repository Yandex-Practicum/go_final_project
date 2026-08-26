package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

// afterNow проверяет, находится ли дата строго после указанного момента.
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// NextDate вычисляет следующую дату выполнения задачи.
// согласно переданному правилу повторения.
//
// dstart должен указывать дату в формате 20060102.
// Итоговый результат всегда должен быть позже now.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty")
	}

	date, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("parse start date: %w", err)
	}

	// Разбиваем правила повторения на отдельные части.
	// Например: "d 7" превращается в ["d", "7"]
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("invalid repeat rule")
	}

	switch parts[0] {

	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid day repeat format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("parse day interval: %w", err)
		}

		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("invalid day interval: %d", interval)
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateLayout), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly repeat format")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateLayout), nil

	case "w":
		// Массив используется как набор допустимых дней недели.
		// 1 - понедельник, 7 - воскресенье.
		var weekday [8]bool

		if len(parts) != 2 {
			return "", errors.New("invalid weekday repeat format")
		}
		days := strings.Split(parts[1], ",")

		for _, day := range days {
			dayNum, err := strconv.Atoi(day)
			if err != nil {
				return "", fmt.Errorf("invalid day format: %w", err)
			}
			if dayNum < 1 || dayNum > 7 {
				return "", fmt.Errorf("invalid weekday: %d", dayNum)
			}
			weekday[dayNum] = true
		}
		for {
			date = date.AddDate(0, 0, 1)

			weekdayNum := int(date.Weekday())

			// time.Weekday() использует воскресенье как 0.
			// поэтому преобразовываем его в 7,
			// чтобы соответствовать формату повторения
			if weekdayNum == 0 {
				weekdayNum = 7
			}

			if weekday[weekdayNum] && afterNow(date, now) {
				return date.Format(dateLayout), nil
			}
		}

	case "m":
		// Индексы массивов соответствуют номерам месяцев и дней:
		// month - 1-12, day - 1-31.
		var month [13]bool
		var day [32]bool

		// Специальные значения -1 и -2 обозначают последний
		// и предпоследний день месяца соответственно.
		var lastDay, prevDay bool

		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid month repeat format")
		}

		days := strings.Split(parts[1], ",")
		for _, d := range days {
			dayNum, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("invalid day format: %w", err)
			}
			if dayNum >= 1 && dayNum <= 31 {
				day[dayNum] = true
			} else if dayNum == -1 {
				lastDay = true
			} else if dayNum == -2 {
				prevDay = true
			} else {
				return "", fmt.Errorf("invalid day: %d", dayNum)
			}
		}

		// Если месяцы не указаны явно, правило повторения действует каждый месяц.
		for i := 1; i <= 12; i++ {
			month[i] = true
		}

		if len(parts) == 3 {
			month = [13]bool{}

			months := strings.Split(parts[2], ",")

			for _, m := range months {
				monthNum, err := strconv.Atoi(m)
				if err != nil {
					return "", fmt.Errorf("invalid month format: %w", err)
				}
				if monthNum < 1 || monthNum > 12 {
					return "", fmt.Errorf("invalid month: %d", monthNum)
				}

				month[monthNum] = true
			}
		}

		limit := now.AddDate(2, 0, 0)

		for {
			date = date.AddDate(0, 0, 1)
			if date.After(limit) {
				return "", errors.New("repeat rule never matches within a reasonable range")
			}

			monthNum := int(date.Month())

			if !month[monthNum] {
				continue
			}

			dayNum := date.Day()

			if day[dayNum] && afterNow(date, now) {
				return date.Format(dateLayout), nil
			}

			// Первый день следующего месяца минус один день - последний день текущего месяца.
			firstNextMonth := time.Date(
				date.Year(),
				date.Month()+1,
				1,
				0, 0, 0, 0,
				date.Location(),
			)
			lastDayOfMonth := firstNextMonth.AddDate(0, 0, -1)

			if lastDay && date.Day() == lastDayOfMonth.Day() && afterNow(date, now) {
				return date.Format(dateLayout), nil
			}
			if prevDay && date.Day() == lastDayOfMonth.Day()-1 && afterNow(date, now) {
				return date.Format(dateLayout), nil
			}

		}

	default:
		return "", errors.New("unsupported repeat rule")
	}

}
