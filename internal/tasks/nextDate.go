package tasks

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lezhnevay/go_final_project_lezhnev/configs"
)

func NextDate(now time.Time, date string, repeat string) (string, error) { // функция возвращает следующую дату
	//проверка на пустой repeat
	if repeat == "" {
		return "", fmt.Errorf("не указан repeat")
	}

	startDate, err := time.Parse(configs.DateFormat, date)
	//проверка на неверный формат даты
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}
	//разделяем правило на тип правила и его аргумент
	ruleSplited := strings.Split(repeat, " ")
	//тип правила
	ruleType := ruleSplited[0]

	switch ruleType {
	case "d":
		if len(ruleSplited) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}
		//количество дней для переноса задачи
		daysToMove, err := strconv.Atoi(ruleSplited[1])

		if err != nil {
			return "", err
		}
		if daysToMove > 400 {

			return "", fmt.Errorf("количество дней не должно превышать 400")
		}
		newDate := startDate.AddDate(0, 0, daysToMove)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем на указанное количество дней
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysToMove)
		}
		return newDate.Format(configs.DateFormat), nil

	case "y":
		newDate := startDate.AddDate(1, 0, 0)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем еще на год
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(configs.DateFormat), nil

	default:
		return "", fmt.Errorf("некорректный тип правила")
	}
}
