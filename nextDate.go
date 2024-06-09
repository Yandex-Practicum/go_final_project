package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		err := errors.New("Not specified repeat rule") //дефотлтная ошибка для вывода при некорректном вводе
		return "", err                                 //пустой repeat
	}

	parsedDate, err := time.Parse(dataFormat, date)
	if err != nil {
		return "", err //неверный формат даты
	}
	commands := strings.Split(repeat, " ") //вычленяем команды для повторов
	ruleType := commands[0]                //берем тип команды

	newDate := parsedDate
	dateFound := false
	fmt.Println(repeat)
	switch ruleType {

	case "d":
		if len(commands) < 2 {
			err := errors.New("Invalid repeat rule for d")
			return "", err
		}
		repeatPeriod, err := strconv.Atoi(commands[1]) //находим дни для переноса задачи

		if err != nil {
			return "", err
		}
		if repeatPeriod > 400 {
			err := errors.New("Invalid repeat rule for d") //дефотлтная ошибка для вывода при некорректном вводе
			return "", err                                 //если агрумент для >400
		}
		newDate = newDate.AddDate(0, 0, repeatPeriod)
		for dateFound == false {
			if (newDate.After(now) && newDate.After(parsedDate)) || now.Format(dataFormat) == newDate.Format(dataFormat) { //если ближайшая дата перевалила за текущий день
				dateFound = true //ну вот тут вопрос как правильнее, так как при return мы все равно выходим из функции, стоит ли поднимать флаг и стоит ли делать break
				return newDate.Format(dataFormat), nil
			}
			newDate = newDate.AddDate(0, 0, repeatPeriod)
		}
		return "", err

	case "y":
		for dateFound == false {
			if newDate.After(now) && newDate.After(parsedDate) { //если ближайшая дата перевалила за текущий день
				dateFound = true

				return newDate.Format(dataFormat), nil
			}
			newDate = newDate.AddDate(1, 0, 0) //находим ближайшую новую дату
		}
		return "", err

	case "w":

		if len(commands) < 2 {
			err := errors.New("Invalid repeat rule for w") //дефолтная ошибка для вывода при некорректном вводе
			return "", err
		}

		daysToRepeatStr := commands[1]                           //находим агрумент для повтора при команде повтора на неделю
		daysToRepeatArray := strings.Split(daysToRepeatStr, ",") //записываем в массив номера дней недели
		var daysToRepeatInts []int                               //записываем в массив номера дней недели

		for _, day := range daysToRepeatArray {

			dayInt, err := strconv.Atoi(day) //проверяем что день недели корректный
			if err != nil {
				err := errors.New("Invalid day for w")
				return "", err
			}
			if dayInt > 7 { // все еще проверяем
				err := errors.New("Invalid day for w") //дефолтная ошибка для вывода при некорректном вводе
				return "", err
			}
			daysToRepeatInts = append(daysToRepeatInts, dayInt)
		}

		sort.Ints(daysToRepeatInts)
		weekAdd := 0
		for dateFound == false {
			sunday := time.Date(newDate.Year(), time.Month(newDate.Month()), newDate.Day()-int(newDate.Weekday())+(weekAdd*7), 0, 0, 0, 0, time.UTC)
			for _, day := range daysToRepeatInts {

				// Вычисляем конкретный день недели

				newDate := sunday.AddDate(0, 0, day)

				if newDate.After(now) && newDate.After(parsedDate) { //если ближайшая дата перевалила за текущий день
					dateFound = true
					return newDate.Format(dataFormat), nil
				}
			}
			weekAdd += 1
		}
		return "", err

	case "m":
		daysToRepeatStr := commands[1]                           //получаем аргументы для повторяемых дней
		daysToRepeatArray := strings.Split(daysToRepeatStr, ",") //переносим дни в массив
		var daysToRepeatInts []int
		for _, day := range daysToRepeatArray {
			dayInt, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}
			if dayInt == 0 || dayInt > 31 || dayInt < -2 { // все еще проверяем
				err := errors.New("Invalid repeat rule for m")
				return "", err
			}
			daysToRepeatInts = append(daysToRepeatInts, dayInt)
		}

		everyMonthSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		monthSlice := make([]int, 1, 1)

		if len(commands) > 2 { //если есть аргумент для месяца
			monthToRepeatStr := commands[2]                            //получаем агрументы для повторяемых месяцев
			monthToRepeatArray := strings.Split(monthToRepeatStr, ",") // переносим месяцы в массив
			for _, month := range monthToRepeatArray {
				monthInt, err := strconv.Atoi(month) //получаем номер месяца из строки
				if err != nil {
					return "", err
				}
				if monthInt > 12 || monthInt <= 0 {
					err := errors.New("Invalid day for m") //дефолтная ошибка для вывода при некорректном вводе
					return "", err
				}
				monthSlice = append(monthSlice, monthInt)
			}
		} else {
			monthSlice = everyMonthSlice
		}
		sort.Ints(monthSlice)
		sort.Ints(daysToRepeatInts)

		year := newDate.Year()
		for dateFound == false {
			for _, month := range monthSlice {
				for _, day := range daysToRepeatInts {
					if day == 0 || day < -2 {
						err := errors.New("Invalid day for m") //дефолтная ошибка для вывода при некорректном вводе
						return "", err
					}
					newDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

					if day < 0 && day > -3 {
						newDate = time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.UTC)
					} else if day != newDate.Day() {
						newDate = time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
					}
					if newDate.After(now) && newDate.After(parsedDate) { //если ближайшая дата перевалила за текущий день
						dateFound = true
						return newDate.Format(dataFormat), err
					}
				}
			}
			year++
		}

		return "", err

	default:
		fmt.Println("Err", repeat)
		err := errors.New("Invalid repeat rule") //дефотлтная ошибка для вывода при некорректном вводе
		return "", err
	}
}

func nextDateHandler(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(dataFormat, now)
	if err != nil {
		return
	}
	response, err := NextDate(nowTime, date, repeat)
	io.WriteString(res, response)
}
