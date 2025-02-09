package dategen

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Функция возвращает следующую дату повторения задачи и ошибку
// Возвращаемая дата должна быть больше даты, указанной в переменной now
//
// Параметры:
//
// now    — время от которого ищется ближайшая дата;
// date   — исходное время в формате 20060102, от которого начинается отсчёт повторений;
// repeat — правило повторения в описанном выше формате.
func NextDate(now time.Time, date string, repeat string) (string, error) {

	repeatSl := strings.Split(repeat, " ")

	// Проверка корректности содержимого date
	//
	dateP, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("Неверный формат данных date: " + date)
	}

	// Проверки корректности содержимого repeat
	// и формирование новых дат
	if (len(repeatSl) <= 1) && (repeat != "y") {
		return "", errors.New("Неверный формат данных 001 repeat: " + repeat)
	}

	switch repeatSl[0] {

	case "d":
		v, err := strconv.Atoi(repeatSl[1])

		// Проверка
		//
		if err != nil {
			return "", errors.New("Неверный формат данных 002 repeat: " + repeat)
		}
		if v > 400 {
			return "", errors.New("Неверный формат данных 003 repeat: " + repeat)
		}
		if v < 1 {
			return "", errors.New("Неверный формат данных 004 repeat: " + repeat)
		}

		// Вычисление новой даты
		//
		if dateP.Before(now) {
			for dateP.Before(now) {
				dateP = dateP.AddDate(0, 0, v)
			}
		} else {
			dateP = dateP.AddDate(0, 0, v)
		}

		//for dateP.Before(now) {
		//	dateP = dateP.AddDate(0, 0, v)
		//}

		return dateP.Format("20060102"), nil // возврат результата по d

	case "y":
		// Проверка
		//
		if len(repeatSl) > 1 {
			return "", errors.New("Неверный формат данных 005 repeat: " + repeat)
		}
		//================================
		// Вычисление новой даты

		if dateP.After(now) {
			dateP = dateP.AddDate(1, 0, 0)
		}

		for dateP.Before(now) {
			dateP = dateP.AddDate(1, 0, 0)
		}

		//================================

		return dateP.Format("20060102"), nil // возврат результата по y

	case "w":
		daysSl := strings.Split(repeatSl[1], ",")

		// Проверка
		//
		if len(daysSl) > 7 {
			return "", errors.New("Неверный формат данных 006 repeat: " + repeat)
		}

		var nd int

		for _, v := range daysSl {
			nd, err = strconv.Atoi(v)
			if err != nil {
				return "", errors.New("Неверный формат данных 007 repeat: " + repeat)
			}
			if nd > 7 {
				return "", errors.New("Неверный формат данных 008 repeat: " + repeat)
			}
			if nd < 1 {
				return "", errors.New("Неверный формат данных 009 repeat: " + repeat)
			}
		}

		// Вычисление новой даты
		//
		var dayOfWeek = map[int]string{
			1: "Monday",
			2: "Tuesday",
			3: "Wednesday",
			4: "Thursday",
			5: "Friday",
			6: "Saturday",
			7: "Sunday",
		}

		var en = true
		for en {
			now = now.AddDate(0, 0, 1)
			if dayOfWeek[nd] == now.Weekday().String() {
				en = false
			}
		}

		//return now.Format("20060102"), nil // возврат результата по w

		return "", errors.New("неподдерживаемый формат")

	case "m":
		switch len(repeatSl) {

		case 2:
			daysSl := strings.Split(repeatSl[1], ",")

			for _, v := range daysSl {
				i, err := strconv.Atoi(v)
				if err != nil {
					return "", errors.New("Неверный формат данных 010 repeat: " + repeat)
				}
				if i > 31 {
					return "", errors.New("Неверный формат данных 011 repeat: " + repeat)
				}
				if i < -31 {
					return "", errors.New("Неверный формат данных 012 repeat: " + repeat)
				}
			}

		case 3:
			daysSl := strings.Split(repeatSl[2], ",")

			for _, v := range daysSl {
				i, err := strconv.Atoi(v)
				if err != nil {
					return "", errors.New("Неверный формат данных 013 repeat: " + repeat)
				}
				if i > 12 {
					return "", errors.New("Неверный формат данных 014 repeat: " + repeat)
				}
				if i < 1 {
					return "", errors.New("Неверный формат данных 015 repeat: " + repeat)
				}
			}

		default:
			return "", errors.New("Неверный формат данных 016 repeat: " + repeat)

		}

	default:
		return "", errors.New("Неверный формат данных 017 repeat: " + repeat)

	}

	return "", errors.New("неподдерживаемый формат")
}
