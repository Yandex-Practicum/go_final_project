package main

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// Map'a дней недели
var week = map[string]string{
	"Monday":    "1",
	"Tuesday":   "2",
	"Wednesday": "3",
	"Thursday":  "4",
	"Friday":    "5",
	"Saturday":  "6",
	"Sunday":    "7",
}

// Функция возвращающая првило и слайс аргументов
func parseRule(repeat string) (string, []string) {
	var repeatRule string
	slice := strings.Split(repeat, " ")
	repeatRule = slice[0]
	return repeatRule, slice[1:]
}

// Функция поиска ближайжей к target даты с учетом числа заданного правилом
func desiredDateWithPositives(day int, target time.Time) time.Time {
	for {
		target = target.AddDate(0, 0, 1)
		if target.Day() == day {
			break
		}
	}
	return target
}

// Функция возвращающая последний день месяца
func lastDayOfMonth(month time.Month, target time.Time) time.Time {
	target = time.Date(target.Year(), month, 1, 0, 0, 0, 0, target.Location())
	return target.AddDate(0, 1, -1)
}

// Функция поиска ближайжей к now даты с учетом чисел и месяцев заданных правилом
func desiredDateWithMonth(daysSlice []int, monthsSlice []int, target time.Time) time.Time {
	datesMap := make(map[time.Time]struct{})
	for _, month := range monthsSlice {
		for _, day := range daysSlice {
			if day < 0 {
				// если значение правила -1 или -2 day = последнийдень заданного месяца + (значение првила +1)
				day = lastDayOfMonth(time.Month(month), target).AddDate(0, 0, day+1).Day()
			}
			// задаем дату с нужным месяцем и днем, год равен году target
			date := time.Date(target.Year(), time.Month(month), day, 0, 0, 0, 0, target.Location())
			if date.Before(target) { // если дата раньше target -> добавляем год
				date = date.AddDate(1, 0, 0)
			}
			datesMap[date] = struct{}{}
		}
	}
	for {
		target = target.AddDate(0, 0, 1)
		if _, ok := datesMap[target]; ok {
			break
		}
	}
	return target
}

// Получение следующей даты
func NextDate(now time.Time, date string, repeat string) (string, error) {
	var result string
	target, err := time.Parse(format, date)
	if err != nil {
		return "", err
	}
	rule, argSlice := parseRule(repeat)

	switch rule {
	//Правило "d"
	case "d":
		if len(argSlice) < 1 { //проверка наличия аргумента
			return "", ErrWrongFormat
		}
		days, err := strconv.Atoi(argSlice[0])
		if err != nil {
			return "", err
		}

		if days > 400 { //проверка на превышение лимита
			return "", ErrExceededLimit
		}
		if target.Before(now) {
			for target.Before(now) {
				target = target.AddDate(0, 0, days)
			}
		} else {
			target = target.AddDate(0, 0, days)
		}
		result = target.Format(format)
	//Правило "y"
	case "y":
		if len(argSlice) > 0 {
			return "", ErrWrongFormat
		}
		if target.Before(now) {
			for target.Before(now) {
				target = target.AddDate(1, 0, 0)
			}
		} else {
			target = target.AddDate(1, 0, 0)
		}
		result = target.Format(format)
	//Правило "w"
	case "w":
		if len(argSlice) < 1 { //проверка наличия аргументов
			return "", ErrWrongFormat
		}
		weekDaysSlice := strings.Split(argSlice[0], ",") //слайс аргументов
		for _, v := range weekDaysSlice {                //проверка значений дней недели
			if v < "1" || v > "7" {
				return "", ErrInvalidValue
			}
		}
		//перевод слайса в map'у
		weekDays := make(map[string]struct{})
		for _, v := range weekDaysSlice {
			weekDays[v] = struct{}{}
		}
		target = now.AddDate(0, 0, 1) //перевод даты на следующий день относительно now

		//ищем совпадения с нужными днями недели
		for {
			if _, ok := weekDays[week[target.Weekday().String()]]; ok { // если совпадение найдено -> выходим из цикла
				break
			}
			target = target.AddDate(0, 0, 1) // прибавляем день
		}
		result = target.Format(format)
	//Правило "m"
	case "m":
		if len(argSlice) < 1 || len(argSlice) > 2 { //проверка на верное кол-во числовых последовательностей
			return "", ErrWrongFormat
		}
		daysStringSlice := strings.Split(argSlice[0], ",") //создание слайса дней из первой последовательности
		daysSlice := make([]int, 0, 31)
		for _, v := range daysStringSlice { //перевод из string в int
			day, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if day < -2 || day > 31 || day == 0 {
				return "", ErrInvalidValue
			}
			daysSlice = append(daysSlice, day)
		}
		if target.Before(now) { //передаем в функцию нужную дату
			target = now
		}
		switch len(argSlice) {
		//если последовательность чисел только одна (дни)
		case 1:
			dates := make([]time.Time, 0, 31)
			for _, day := range daysSlice {
				if day > 0 {
					dates = append(dates, desiredDateWithPositives(day, target))
				} else {
					dates = append(dates, lastDayOfMonth(target.Month(), target).AddDate(0, 0, day+1))
				}
			}
			sort.Slice(dates, func(i, j int) bool { //сортировка слайса
				return dates[i].Before(dates[j])
			})
			result = dates[0].Format(format) //первая дата будет ближайшей к target и удовлетворять правилу
		//если последовательностей чисел две (дни + месяцы)
		case 2:
			monthsStringSlice := strings.Split(argSlice[1], ",")
			monthsSlice := make([]int, 0, 12)
			for _, v := range monthsStringSlice {
				month, err := strconv.Atoi(v)
				if err != nil {
					return "", err
				}
				if month < 1 || month > 12 { // проверяем верность переданных значений
					return "", ErrInvalidValue
				}
				monthsSlice = append(monthsSlice, month) //заполняем слайс
			}
			result = desiredDateWithMonth(daysSlice, monthsSlice, target).Format(format)
		}
	// пустое правило
	case "":
		return "", nil
	//неверное правило
	default:
		return "", ErrWrongFormat
	}
	return result, nil
}

// Проверка и исправление, если это возможно, переданного задания
func checkTask(task *Task) (err error) {
	if task.Title == "" { // проверка наличия заголовка
		return ErrNoTitle
	}
	today := time.Now().Format(format)
	if task.Date == "" {
		task.Date = today
	}
	if _, err := time.Parse(format, task.Date); err != nil { //проверка даты
		return err
	}
	if task.Date < today && task.Repeat == "" { //если дата меньше текущей и правило пустое подставляем сегодняшнее число
		task.Date = today
	}
	if task.Repeat != "" && task.Date < today { //если правило пустое и дата меньше текущей используем NextDate()
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	return nil
}
