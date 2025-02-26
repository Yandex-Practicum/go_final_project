package date

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	var nextDateString string
	var nextDate time.Time

	// Проверка даты начала отсчёта
	_, err := Validation(date)
	if err != nil {
		fmt.Println("Неверное значение даты начала отсчета для повторений. ", err)
		return nextDateString, err
	}

	// Проверка корректности интервалов для повторений
	err = RepeatValidation(repeat)
	if err != nil {
		fmt.Println("Неверное значение для повторений. ", err)
		return nextDateString, err
	}

	// Если строка для повторений пуста
	if repeat == "" {
		return nextDateString, nil
	}

	switch string(repeat[0]) {
	case "y":
		nextDate, err = yearDateToRepeat(now, date, repeat)
		if err != nil {
			fmt.Println("Год для повторений не корректен ", err)
			return nextDateString, err
		}
	case "d":
		nextDate, err = dayDateToRepeat(now, date, repeat)
		if err != nil {
			fmt.Println("День для повторений не корректен ", err)
			return nextDateString, err
		}
	case "w":
		nextDate, err = weekDateToRepeat(now, date, repeat)
		if err != nil {
			fmt.Println("Неделя для повторений не корректна ", err)
			return nextDateString, err
		}

	default:
		err0 := errors.New("значение для повторений не корректно")
		return nextDateString, err0
	}

	nextDateString = fmt.Sprint(nextDate.Format(FormatDate))
	return nextDateString, nil
}

// yearDateToRepeat - возвращает следующую дату для ежегодного повторения
func yearDateToRepeat(now time.Time, date string, repeat string) (time.Time, error) {

	var nextDate time.Time

	// На случай, если формат ежегодных повторений будет меняться

	if len(repeat) != 1 {
		err0 := errors.New("такие данные пока не обрабатываются")
		return nextDate, err0
	}

	dateStart, err := time.Parse(FormatDate, date)
	if err != nil {
		fmt.Println("Не верные данные date. ", err)
		return nextDate, err
	}

	nextDate = dateStart.AddDate(1, 0, 0)
	for now.After(nextDate) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}

	return nextDate, nil
}

// dayDateToRepeat - возвращает следующую дату для повторения через несколько дней
func dayDateToRepeat(now time.Time, date string, repeat string) (time.Time, error) {

	var nextDate time.Time

	dateStart, err := time.Parse(FormatDate, date)
	if err != nil {
		fmt.Println("Не верные данные date. ", err)
		return nextDate, err
	}

	dayData := strings.Split(repeat, " ")

	dayCount, err := strconv.Atoi(dayData[1])
	if err != nil {
		fmt.Println("Неверный формат дней для повторений. ", err)
		return nextDate, err
	}

	nextDate = dateStart.AddDate(0, 0, dayCount)
	for now.After(nextDate) {
		nextDate = nextDate.AddDate(0, 0, dayCount)
	}

	return nextDate, nil
}

// weekDateToRepeat - возвращает следующую дату для еженедельного повторения
func weekDateToRepeat(now time.Time, date string, repeat string) (time.Time, error) {

	week := map[time.Weekday]int{time.Monday: 1, time.Tuesday: 2, time.Wednesday: 3, time.Thursday: 4, time.Friday: 5, time.Saturday: 6, time.Sunday: 7}

	var nextDate time.Time

	dateStart, err := time.Parse(FormatDate, date)
	if err != nil {
		fmt.Println("Не верные данные date. ", err)
		return nextDate, err
	}

	weekDays := strings.Split(repeat, " ")

	var weekDaysInt []int
	for _, value := range weekDays {
		dayNumber, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Не могу конвертировать в цифру значение недельного повторения ", err)
			return nextDate, err
		}
		weekDaysInt = append(weekDaysInt, dayNumber)
	}

	// Прибавляем с шагом в неделю пока не найдём.
	for i := 0; i < 56; i++ {
		dateStart = dateStart.AddDate(0, 0, 1*i)
		for _, value := range weekDaysInt {

			d := week[dateStart.Weekday()]
			w := value

			delta := (7 - d + w) % 7
			nextDate = dateStart.AddDate(0, 0, delta)

			if nextDate.After(now) {
				break
			}
		}
	}

	// На случай если что-то пошло не так
	if !nextDate.After(now) {
		err0 := errors.New("не получилось найти следующую дату недельного повторения")
		fmt.Println(err0)
		return nextDate, err0
	}

	return nextDate, nil
}
