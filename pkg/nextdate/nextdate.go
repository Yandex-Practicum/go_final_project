package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "20060102"

// если переданная дата больше текущей - true
func afterNow(date, now time.Time) bool {
	// обнуляем время
	date = date.Truncate(24 * time.Hour)
	now = now.Truncate(24 * time.Hour)
	return date.After(now)
}

// меняет воскресенье с 0 на 7
func weekdayToISO(w time.Weekday) int {
	if w == time.Sunday {
		return 7
	}
	return int(w)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// парсим дату
	dateStart, err := time.Parse(TimeFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга даты: %w", err)
	}

	if len(repeat) == 0 {
		return "", fmt.Errorf("отсутствует правило")
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "y":
		for {
			// увеличиваем год на 1
			dateStart = dateStart.AddDate(1, 0, 0)

			// если дата равна или больше текущей, выходим
			if dateStart.Year() >= now.Year() {
				return dateStart.Format(TimeFormat), nil
			}
		}

	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("нужен аргумент после d")
		}

		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("ошибка парсинга дней: %w", err)
		}

		if num < 1 || num > 400 {
			return "", fmt.Errorf("некорректный день")
		}

		for {

			dateStart = dateStart.AddDate(0, 0, num)
			if afterNow(dateStart, now) {
				break
			}
		}
		return dateStart.Format(TimeFormat), nil

	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("нужен аргумент после w")
		}

		days := []int{}
		numbers := strings.Split(parts[1], ",")
		for _, number := range numbers {
			day, err := strconv.Atoi(number)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("ошибка парсинга дня недели: %w", err)
			}
			days = append(days, day)
		}

		// ищем ближайший день недели
		for {
			if afterNow(dateStart, now) {
				currentDay := weekdayToISO(dateStart.Weekday())
				for _, d := range days {
					if d == currentDay {
						return dateStart.Format(TimeFormat), nil
					}
				}
			}
			dateStart = dateStart.AddDate(0, 0, 1)
		}

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("нужен аргумент после m")
		}

		// парсим дни
		dayParts := strings.Split(parts[1], ",")
		days := []int{}
		for _, el := range dayParts {
			day, err := strconv.Atoi(el)
			if err != nil || day == 0 || day < -31 || day > 31 {
				return "", fmt.Errorf("неверное значение дня")
			}
			days = append(days, day)
		}

		// парсим месяцы
		months := []int{}
		if len(parts) >= 3 {
			monthParts := strings.Split(parts[2], ",")
			for _, el := range monthParts {
				month, err := strconv.Atoi(el)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("неверное значение месяца")
				}
				months = append(months, month)
			}
		}

		for {
			if afterNow(dateStart, now) {
				year, month := dateStart.Year(), dateStart.Month()
				currentDay := dateStart.Day()

				// учитываем месяцы
				if len(months) > 0 {
					matchMonth := false
					for _, m := range months {
						if int(month) == m {
							matchMonth = true
							break
						}
					}
					if !matchMonth {
						dateStart = dateStart.AddDate(0, 0, 1)
						continue
					}
				}

				// учитываем дни
				lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
				for _, d := range days {
					if d == -3 {
						return "", fmt.Errorf("неверное значение")
					}
					targetDay := d
					if targetDay < 0 {
						targetDay = lastDayOfMonth + d + 1
					}
					if targetDay == currentDay {
						return dateStart.Format(TimeFormat), nil
					}
				}
			}
			dateStart = dateStart.AddDate(0, 0, 1)
		}

	default:
		return "", fmt.Errorf("неизвестное правило")
	}
}
