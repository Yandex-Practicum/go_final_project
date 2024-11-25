package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go_final_project/internal/constants"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	//Проверка на наличие правила повторения
	if repeat == "" {
		return "", errors.New("отсутствует правило повторения")
	}

	//Проверка формата даты
	taskDate, err := time.Parse(constants.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат даты: %v", err)
	}

	//Правило повторения при "d 1"
	if repeat == "d 1" && !taskDate.After(now) {
		return now.Format(constants.DateFormat), nil
	}

	//проверка соответствия условиям правила повторения
	rpt := strings.Split(repeat, " ")
	for {
		if rpt[0] == "y" && len(rpt) < 2 {
			taskDate = taskDate.AddDate(1, 0, 0)
		} else if rpt[0] == "d" && len(rpt) == 2 {
			days, err := strconv.Atoi(rpt[1])
			if err != nil || days < 1 || days > 400 {
				return "", fmt.Errorf("некорректное число дней: d %s", rpt[1])
			}
			taskDate = taskDate.AddDate(0, 0, days)
		} else {
			return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
		}

		if taskDate.After(now) {
			return taskDate.Format(constants.DateFormat), nil
		}
	}

}
