package nextdate

import (
	"errors"
	"strconv"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if len(date) == 0 || len(repeat) == 0 {
		return "", errors.New("все параметры функции обязательны")
	}

	timeDate, err := time.Parse("20060102", date)
	//если не удалось разобрать date
	if err != nil {
		return "", err
	}

	var handler func()

	invalidRepeatErr := errors.New("не удалось разобрать repeat")

	if repeat[0] == 'y' {
		handler = func() { timeDate = timeDate.AddDate(1, 0, 0) }
	} else if repeat[0] == 'd' {

		if len(repeat) <= 2 || repeat[1] != ' ' {
			return "", invalidRepeatErr
		}

		if repeat == "d 1" && !timeDate.After(now) {
			return now.Format("20060102"), nil
		}

		i, err := strconv.ParseInt(repeat[2:], 10, 32)
		if err != nil || i > 400 {
			return "", invalidRepeatErr
		}

		handler = func() { timeDate = timeDate.AddDate(0, 0, int(i)) }

	} else {
		return "", invalidRepeatErr
	}

	for {
		handler()
		if timeDate.After(now) {
			break
		}
	}

	return timeDate.Format("20060102"), nil //допилить до строки
}
