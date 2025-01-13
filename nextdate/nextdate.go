package nextdate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const ParseDate = "20060102"

func CalcNextDate(now, date, repeat string) (string, error) {
	// получаю правила повторения задач
	rule, err := ParseRepeatRules(repeat)
	if err != nil {
		return "", errors.New("Формат правила повторения не соблюден")
	}
	// парсинг полученных дат
	nowTime, dateTime, err := ParsingDates(now, date)
	if err != nil {
		return "", errors.New("Некорректный формат даты")
	}
	//вычисление дня переноса задачи
	if rule[0] == "d" {
		resultDate, err := CountDateRepeatDay(rule, nowTime, dateTime)
		if err != nil {
			return "", errors.New("Формат правила повторения не соблюден")
		}
		return resultDate, nil
	} else {
		resultDate, err := CountDateRepeatRule(nowTime, dateTime)
		if err != nil {
			return "", errors.New("Формат правила повторения не соблюден")
		}
		return resultDate, nil
	}
}

func ParseRepeatRules(rule string) ([]string, error) {
	repeatRule := strings.Split(rule, " ")
	if (repeatRule[0] == "d" && len(repeatRule) == 2) || (repeatRule[0] == "y" && len(repeatRule) == 1) {
		return repeatRule, nil
	} else {
		return repeatRule, errors.New("Формат правила повторения не соблюден")
	}
}

func ParsingDates(now, date string) (time.Time, time.Time, error) {
	nowTime, err := time.Parse(ParseDate, now)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Некорректный формат даты")
	}
	dateTime, err := time.Parse(ParseDate, date)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("Некорректный формат даты")
	}
	return nowTime, dateTime, nil
}

func CountDateRepeatDay(rules []string, nowTime, dateTime time.Time) (string, error) {
	subtraction := dateTime.Sub(nowTime)
	days, err := strconv.Atoi(rules[1])
	if (err != nil) || (days > 400) {
		return "", errors.New("Формат правила повторения не соблюден")
	}
	if int(subtraction.Hours()) > 0 {
		dateTime = dateTime.AddDate(0, 0, days)
		return dateTime.Format(ParseDate), nil
	}
	for int(subtraction.Hours()) <= 0 {
		dateTime = dateTime.AddDate(0, 0, days)
		subtraction += time.Duration(days * 24 * int(time.Hour))
	}
	return dateTime.Format(ParseDate), nil
}

func CountDateRepeatRule(nowTime, dateTime time.Time) (string, error) {
	//определяем високосный год или нет
	ageStringdate := dateTime.Format(ParseDate)
	ageStringnow := nowTime.Format(ParseDate)
	resDate, _ := strconv.Atoi(ageStringdate)
	resNow, _ := strconv.Atoi(ageStringnow)
	ageDate := int(resDate) / int(10000)
	ageNow := int(resNow) / int(10000)

	if ageDate >= ageNow {
		dateTime = dateTime.AddDate(1, 0, 0)
		return dateTime.Format(ParseDate), nil
	}
	for ageDate < ageNow {
		dateTime = dateTime.AddDate(1, 0, 0)
		if (ageDate%4 == 0) && (ageDate%100 == 0) && (ageDate%400 == 0) {
			ageDate += 1
		} else {
			ageDate += 1
		}
	}
	return dateTime.Format(ParseDate), nil
}
