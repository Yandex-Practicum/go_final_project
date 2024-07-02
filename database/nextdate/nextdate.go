package nextdate

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// переменные используемые всеми функциями в этом package
var (
	startDate  time.Time
	now        time.Time
	nextDate   string
	dateFormat string
)

// NextDate возвращает дату и ошибку, исходя из правил указанных в repeat.
func NextDate(nowArg time.Time, date string, repeat string) (string, error) {
	// Если dateFormat ещё не инициализирована, берём её из .env файла
	if len(dateFormat) == 0 {
		dateFormat = os.Getenv("TODO_DATEFORMAT")
	}

	var err error

	if repeat == "" {
		return "", fmt.Errorf("пустая строка в repeat")
	}

	now = nowArg
	startDate, err = time.Parse(dateFormat, date)
	if err != nil {
		return "", err
	}

	repeat = strings.ToLower(repeat)
	prefix := string(repeat[0])
	code, _ := strings.CutPrefix(repeat, prefix)
	code = strings.TrimSpace(code)

	switch prefix {
	case "d":
		_, err = calcD(code)
	case "y":
		calcY()
	case "w":
		_, err = calcW(code)
	case "m":
		_, err = calcM(code)
	default:
		return "", fmt.Errorf("некорректный формат repeat")
	}
	if err != nil {
		return "", err
	}

	return nextDate, nil

}

// listAtoi конвертирует слайс string в слайс int
func listAtoi(list []string) ([]int, error) {
	var resList []int
	for _, str := range list {
		num, err := strconv.Atoi(str)
		if err != nil {
			return []int{}, err
		}
		resList = append(resList, num)
	}
	return resList, nil
}

// getClosestIdx возвращает индекс блжиайшей переменной в слайсе list, превышающей значение target. Если такой перменной нет, возвращает -1.
func getClosestIdx(target int, list []int) int {
	idx := slices.IndexFunc(list, func(d int) bool { return d > target })
	return idx
}

// pickBigger возвращает большее из двух чисел
func pickBigger(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}
