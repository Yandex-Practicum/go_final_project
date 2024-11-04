package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go_final_project/internal/config"
)

// NextDate возвращает следующую дату выполнения задачи на основе правил повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Проверка наличия правила повторения
	if repeat == "" {
		return "", fmt.Errorf("правило повторения обязательно")
	}

	// Парсинг даты задачи
	dateTask, err := time.Parse(config.TimeFormat, date)
	if err != nil {
		return "", err
	}

	// Специальное условие: если правило повторения "d N" и дата совпадает с текущей
	repeatParts := strings.Split(repeat, " ")
	if strings.ToLower(repeatParts[0]) == "d" && dateTask.Equal(now.Truncate(24*time.Hour)) {
		return now.AddDate(0, 0, parseValueB(repeatParts[01])).Format(config.TimeFormat), nil

	}

	// Обработка правил повторения
	switch strings.ToLower(repeatParts[0]) {
	case "d": // Повторение по дням
		if len(repeatParts) < 2 {
			return "", fmt.Errorf("для дней необходимо указать число")
		}

		// Парсинг значения дней
		days, err := parseValue(repeatParts[1])
		if err != nil {
			return "", err
		}
		dateTask = addDate(now, dateTask, 0, 0, days)

	case "y": // Повторение по годам
		dateTask = addDate(now, dateTask, 1, 0, 0)

	case "w": // Повторение по неделям
		if len(repeatParts) < 2 {
			return "", fmt.Errorf("для недель необходимо указать дни")
		}
		dateTask, err = getDateByWeek(now, dateTask, repeatParts[1])
		if err != nil {
			return "", err
		}

	case "m": // Повторение по месяцам
		if len(repeatParts) < 2 {
			return "", fmt.Errorf("для месяцев необходимо указать дни")
		}
		dateTask, err = getDateByMonth(now, dateTask, repeatParts)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("недопустимое правило")
	}

	// Возвращаем следующую дату в нужном формате
	return dateTask.Format(config.TimeFormat), nil
}

// parseValue парсит значение для дней
func parseValue(num string) (int, error) {
	days, err := strconv.Atoi(num) // Преобразование строки в число
	if err != nil {
		return 0, err
	}

	// Проверка допустимости значения
	if days < 0 || days >= 400 {
		return 0, fmt.Errorf("недопустимое значение %d", days)
	}

	return days, nil
}

// parseValue парсит значение для дней
func parseValueB(num string) int {
	days, err := strconv.Atoi(num) // Преобразование строки в число
	if err != nil {
		return 0
	}

	// Проверка допустимости значения
	if days < 0 || days >= 400 {
		return 0
	}

	return days
}

// addDate добавляет к дате год, месяц и день
func addDate(now time.Time, dateTask time.Time, year int, month int, day int) time.Time {
	dateTask = dateTask.AddDate(year, month, day) // Добавление лет, месяцев и дней

	// Увеличиваем дату, если она меньше текущей
	for dateTask.Before(now) {
		dateTask = dateTask.AddDate(year, month, day)
	}

	return dateTask
}

// getDateByWeek возвращает дату выполнения задачи на основе дней недели
func getDateByWeek(now, dateTask time.Time, daysString string) (time.Time, error) {
	daysSlice := strings.Split(daysString, ",") // Разделение строки на дни
	daysOfWeekMap := make(map[int]bool)

	// Парсинг дней недели
	for _, day := range daysSlice {
		numberOfDay, err := strconv.Atoi(day)
		if err != nil {
			return dateTask, err
		}

		if numberOfDay < 1 || numberOfDay > 7 {
			return dateTask, fmt.Errorf("недопустимое значение %d для дня недели", numberOfDay)
		}

		if numberOfDay == 7 {
			numberOfDay = 0 // Преобразование 7 в 0 для соответствия с Go
		}
		daysOfWeekMap[numberOfDay] = true
	}

	// Поиск следующей даты, соответствующей дню недели
	for {
		if daysOfWeekMap[int(dateTask.Weekday())] {
			if now.Before(dateTask) {
				break
			}
		}
		dateTask = dateTask.AddDate(0, 0, 1) // Переход к следующему дню
	}

	return dateTask, nil
}

// getDateByMonth возвращает дату выполнения задачи на основе дней месяца
func getDateByMonth(now, dateTask time.Time, repeat []string) (time.Time, error) {
	daysString := repeat[1]
	monthsString := ""
	if len(repeat) > 2 {
		monthsString = repeat[2]
	}

	daysSlice := strings.Split(daysString, ",")
	monthsSlice := strings.Split(monthsString, ",")

	// Парсинг дней месяца
	daysMap := make(map[int]bool)
	for _, day := range daysSlice {
		numberOfDay, err := strconv.Atoi(day)
		if err != nil {
			return dateTask, err
		}
		if numberOfDay < -2 || numberOfDay > 31 || numberOfDay == 0 {
			return dateTask, fmt.Errorf("недопустимое значение %d для дня месяца", numberOfDay)
		}
		daysMap[numberOfDay] = true
	}

	// Парсинг месяцев
	monthsMap := make(map[int]bool)
	for _, month := range monthsSlice {
		if month == "" {
			continue
		}
		numberOfMonth, err := strconv.Atoi(month)
		if err != nil {
			return dateTask, err
		}
		if numberOfMonth < 1 || numberOfMonth > 12 {
			return dateTask, fmt.Errorf("недопустимое значение %d для месяца", numberOfMonth)
		}
		monthsMap[numberOfMonth] = true
	}

	// Поиск следующей даты, соответствующей месяцу
	for {
		if len(monthsMap) == 0 {
			break
		}

		if monthsMap[int(dateTask.Month())] {
			if now.Before(dateTask) {
				break
			}
		}
		dateTask = dateTask.AddDate(0, 0, 1) // Переход к следующему дню
	}

	// Проверка на совпадение с днями месяца
	for {
		lastDay := time.Date(dateTask.Year(), dateTask.Month()+1, 0, 0, 0, 0, 0, dateTask.Location()).Day()
		predLastDay := lastDay - 1

		key := dateTask.Day()
		switch {
		case lastDay == dateTask.Day():
			if _, ok := daysMap[-1]; ok {
				key = -1 // последний день месяца
			}
		case predLastDay == dateTask.Day():
			if _, ok := daysMap[-2]; ok {
				key = -2 // предпоследний день месяца
			}
		}

		if daysMap[key] {
			if now.Before(dateTask) {
				break
			}
		}
		dateTask = dateTask.AddDate(0, 0, 1) // Переход к следующему дню
	}

	return dateTask, nil
}
