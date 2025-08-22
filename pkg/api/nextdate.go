package api

import (
	"errors"
	"fmt"
	"time"
)

const isoDate = "20060102"

var (
	ErrBadRepeat      = errors.New("invalid or missing repeat rule")
	ErrInvalidDate    = errors.New("invalid date format, expected YYYYMMDD")
	ErrIntervalTooBig = errors.New("day interval exceeds maximum of 400")
	ErrUnsupported    = errors.New("unsupported repeat rule (only d <N> and y)")
)

// NextDate возвращает следующую дату > now по правилу repeat,
// считая от исходной даты dstart (формат YYYYMMDD).
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	// 1. Парсим dstart
	start, err := time.Parse(isoDate, dstart)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidDate, err)
	}

	// 2. Обработка пустого repeat
	if repeat == "" {
		return "", ErrBadRepeat
	}

	// 3. Определяем вид регулярности
	var stepFunc func(time.Time) (time.Time, error)

	switch {
	case repeat == "y":
		// ежегодно
		stepFunc = func(t time.Time) (time.Time, error) {
			return t.AddDate(1, 0, 0), nil
		}

	case len(repeat) > 2 && repeat[:2] == "d ":
		// d N — через N дней
		var days int
		if _, err := fmt.Sscanf(repeat, "d %d", &days); err != nil {
			return "", ErrBadRepeat
		}
		if days < 1 || days > 400 {
			return "", ErrIntervalTooBig
		}
		stepFunc = func(t time.Time) (time.Time, error) {
			return t.AddDate(0, 0, days), nil
		}

	default:
		// остальные правила пока не поддерживаем
		return "", ErrUnsupported
	}

	// 4. Находим первое значение > now
	candidate := start
	for {
		candidate, err = stepFunc(candidate)
		if err != nil {
			return "", err
		}
		// сравниваем даты без времени
		if candidate.After(now) {
			break
		}
	}

	return candidate.Format(isoDate), nil
}
