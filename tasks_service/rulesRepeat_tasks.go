package tasks_service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid 'now' date format", http.StatusBadRequest)
		log.Printf("Неверный формат даты 'now' в nextDateHandler: %v", err)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Ошибка вычисления следующей даты в nextDateHandler: %v", err)
		return
	}

	log.Println("Следующая дата вычислена успешно:", nextDate)
	fmt.Fprint(w, nextDate)
}

// DateFormat — это постоянный формат для дат, используемый в приложении
const DateFormat = "20060102"

// NextDate вычисляет следующую дату для задачи на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения пустое")
	}

	// Разбираем начальную дату задачи с использованием постоянного формата DateFormat
	taskDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	// Начинаем расчет с taskDate или now, в зависимости от того, что больше
	startDate := now
	if taskDate.After(now) {
		startDate = taskDate
	}

	switch {
	case strings.HasPrefix(repeat, "d "): // Обрабатываем правило ежедневного повторения "d <число>"

		//if repeat == "d" {
		//	repeat = "d 1"
		//}

		daysStr := strings.TrimSpace(repeat[2:])
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неверное правило повторения 'd': %v", err)
		}

		taskDate = taskDate.AddDate(0, 0, days)

		for taskDate.Before(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}

		return taskDate.Format(DateFormat), nil

	case repeat == "y": // Обрабатываем правило ежегодного повторения "y"
		// Логика прямо рассчитывает следующую ежегодную дату без проверки, если taskDate после now
		for !taskDate.After(startDate) {
			year := taskDate.Year() + 1
			month := taskDate.Month()
			day := taskDate.Day()

			// Обработка сценария високосного года
			if month == time.February && day == 29 && !isLeapYear(year) {
				taskDate = time.Date(year, time.March, 1, 0, 0, 0, 0, taskDate.Location())
			} else {
				taskDate = time.Date(year, month, day, 0, 0, 0, 0, taskDate.Location())
			}
		}

		return taskDate.Format(DateFormat), nil

	case strings.HasPrefix(repeat, "w "): // Обрабатываем правило еженедельного повторения "w <дни>"
		daysStr := strings.TrimSpace(repeat[2:])
		days := strings.Split(daysStr, ",")
		if len(days) == 0 {
			return "", fmt.Errorf("неверное правило повторения 'w': дни не указаны")
		}

		var daysOfWeek []int
		for _, dayStr := range days {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("неверное правило повторения 'w': неверный день '%s'", dayStr)
			}
			if day == 7 {
				day = 0
			}
			daysOfWeek = append(daysOfWeek, day)
		}
		sort.Ints(daysOfWeek)

		// Устанавливаем начальную дату startDate на максимум из taskDate и now
		startDate = taskDate
		if now.After(taskDate) {
			startDate = now
		}

		initialDate := taskDate

		// Выполняем цикл, пока день недели не совпадет с одним из daysOfWeek
		for !containsInt(daysOfWeek, int(startDate.Weekday())) || !(startDate.YearDay() > initialDate.YearDay()) {
			startDate = startDate.AddDate(0, 0, 1)
		}
		return startDate.Format(DateFormat), nil

	case strings.HasPrefix(repeat, "m "): // Обрабатываем правило ежемесячного повторения "m <дни> [<месяцы>]"
		parts := strings.Split(strings.TrimSpace(repeat[2:]), " ")

		if len(parts) == 0 {
			return "", fmt.Errorf("неверное правило повторения 'm': дни не указаны")
		}

		dayParts := strings.Split(parts[0], ",")
		var daysOfMonth []int
		for _, dayStr := range dayParts {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day == 0 || day < -2 || day > 31 {
				return "", fmt.Errorf("неверное правило повторения 'm': неверный день '%s'", dayStr)
			}
			daysOfMonth = append(daysOfMonth, day)
		}

		var months []int
		if len(parts) > 1 {
			monthParts := strings.Split(parts[1], ",")
			for _, monthStr := range monthParts {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("неверное правило повторения 'm': неверный месяц '%s'", monthStr)
				}
				months = append(months, month)
			}
		}

		sort.Ints(daysOfMonth)
		sort.Ints(months)

		for {
			currentYear, currentMonth := taskDate.Year(), taskDate.Month()

			if len(months) > 0 && !containsInt(months, int(currentMonth)) {
				nextMonth := findNextMonth(int(currentMonth), months)
				if nextMonth <= int(currentMonth) {
					taskDate = time.Date(currentYear+1, time.Month(nextMonth), 1, 0, 0, 0, 0, taskDate.Location())
				} else {
					taskDate = time.Date(currentYear, time.Month(nextMonth), 1, 0, 0, 0, 0, taskDate.Location())
				}
				continue
			}

			// Находим ближайшую допустимую дату в текущем месяце
			var nextValidDate time.Time
			for _, day := range daysOfMonth {
				var candidateDate time.Time
				lastDay := lastDayOfMonth(currentYear, currentMonth)

				if day > 0 {
					if day <= lastDay {
						candidateDate = time.Date(currentYear, currentMonth, day, 0, 0, 0, 0, taskDate.Location())
					} else {
						continue
					}
				} else {
					if -day <= lastDay {
						candidateDate = time.Date(currentYear, currentMonth, lastDay+day+1, 0, 0, 0, 0, taskDate.Location())
					} else {
						continue
					}
				}

				// Находим ближайшую дату после startDate
				if candidateDate.After(startDate) && (nextValidDate.IsZero() || candidateDate.Before(nextValidDate)) {
					nextValidDate = candidateDate
				}
			}

			// Если найдена допустимая дата, возвращаем её
			if !nextValidDate.IsZero() {
				return nextValidDate.Format(DateFormat), nil
			}

			// Если не найдено допустимой даты в текущем месяце, переходим к следующему месяцу
			taskDate = taskDate.AddDate(0, 1, 0)
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), 1, 0, 0, 0, 0, taskDate.Location())
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: '%s'", repeat)
	}
}
