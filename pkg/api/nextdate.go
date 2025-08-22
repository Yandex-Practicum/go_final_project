package api

import (
	"errors"

	"strconv"
	"time"
)

//const DateFormat = "20060102"

func NextDate(nowtime time.Time, dstart string, repeat string) (string, error) {
	//проверка пустой или нет repeat
	if len(repeat) == 0 {

		return "", errors.New("repeat is empty")
	}

	//разбираем dstart на соответствие формата
	parsedTime, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("error parse dstart")
	}

	//если другой формат repeat. в моём случае только 2 варианта "d " и "y"
	if len(repeat) > 3 {
		if repeat[:2] != "d " {
			//fmt.Println("error format repeat")
			return "", errors.New("error format repeat")
		}

	}

	if len(repeat) == 1 {
		if repeat != "y" {
			//fmt.Println("error format repeat")
			return "", errors.New("error format repeat")
		}
	}

	//если repeat == y , т.е. прибавить год
	if repeat == "y" {
		for {
			parsedTime = parsedTime.AddDate(1, 0, 0)
			if parsedTime.After(nowtime) {
				break
			}

		}
		return parsedTime.Format(DateFormat), nil
	}

	repeat = repeat[2:] // Удаляем первые 2 символа
	//если repeat цифра
	num, err := strconv.Atoi(repeat)
	if err != nil {
		return "", errors.New("wrong format repeat")
	}

	if num > 400 {
		return "", errors.New("repeat > 400")
	}

	for {
		parsedTime = parsedTime.AddDate(0, 0, num)
		if parsedTime.After(nowtime) {
			break
		}
	}

	return parsedTime.Format(DateFormat), nil

}
