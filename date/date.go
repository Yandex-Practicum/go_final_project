package date

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// выносим формат даты в константу
const dateFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// проверка, что правило указано
	if len(repeat) == 0 {
		return "", errors.New("правило повторения не указано")
	}

	// проверка, что в параметр date передали корректную дату
	tDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", errors.New("время в переменной date не может быть преобразовано в корректную дату")
	}

	// получаем слайс с подстроками правила повторения
	slRepeat := strings.Fields(repeat)

	// проверка, что правило корректно по количеству элементов
	if len(slRepeat) < 1 || len(slRepeat) > 2 {
		return "", errors.New("правило указано некорректно")
	}

	// инициализируем количество дней для правил
	daysCount := 0

	switch slRepeat[0] {
	case "d":
		if len(slRepeat) != 2 {
			return "", errors.New("правило с днями указано некорректно")
		}

		daysCount, err = strconv.Atoi(slRepeat[1])
		if err != nil || daysCount < 1 || daysCount > 400 {
			return "", errors.New("правило с днями указано некорректно или недопустимое количество дней")
		}
	case "y":
		if len(slRepeat) != 1 {
			return "", errors.New("правило с годами указано некорректно")
		}
	default:
		return "", errors.New("указанный формат правила повторения не поддерживается")
	}

	// считаем для дней
	if daysCount > 0 {

		// применяем правило обязательно хотя бы 1 раз
		tDate = tDate.AddDate(0, 0, daysCount)

		// применяем правило в цикле, пока дата не станет больше now
		for tDate.Format(dateFormat) <= now.Format(dateFormat) {
			tDate = tDate.AddDate(0, 0, daysCount)
		}
	}

	// считаем для лет
	if daysCount == 0 {

		// применяем правило обязательно хотя бы 1 раз
		tDate = tDate.AddDate(1, 0, 0)

		// применяем правило в цикле, пока дата не станет больше now
		for tDate.Format(dateFormat) <= now.Format(dateFormat) {
			tDate = tDate.AddDate(1, 0, 0)
		}
	}

	return tDate.Format(dateFormat), nil
}
