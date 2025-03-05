package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go_final_project/date"
)

// TaskDataValidation - проверяет корректность входящих данных
func TaskDataValidation(httpData []byte) (Task, error) {

	var taskData Task
	now := time.Now()

	// Приводим текущее время к формату
	nowString := now.Format(date.FormatDate)
	now, err := time.Parse(date.FormatDate, nowString)
	if err != nil {
		fmt.Println("Строковые данные даты не корректны. ", err)
		return taskData, err
	}

	// Проверка распаковки JSON
	err = json.Unmarshal(httpData, &taskData)
	if err != nil {
		fmt.Println("не удачно распаковался JSON запрос ", err)
		return taskData, err
	}

	// Проверка id
	if taskData.ID != "" {
		_, err := IDValidation(taskData.ID)
		if err != nil {
			fmt.Println("не верый формат id ", err)
			return taskData, err
		}
	}

	// Проверка заголовка
	if taskData.Title == "" {
		err0 := errors.New("заголовок задачи отсутствует")
		return taskData, err0
	}

	// Проверка значений повторений
	if date.RepeatValidation(taskData.Repeat) != nil {
		return taskData, date.RepeatValidation(taskData.Repeat)
	}

	// Проверка на пустое значение даты
	if taskData.Date == "" {
		taskData.Date = fmt.Sprint(now.Format(date.FormatDate))
	}

	// Проверка на корректность даты
	dateToTask, err := date.Validation(taskData.Date)
	if err != nil {
		fmt.Println(taskData.Date, err)
		return taskData, err
	}
	taskData.Date = fmt.Sprint(dateToTask.Format(date.FormatDate))

	// Если дата меньше сегодняшей
	// Значение переданной даты. Ошибку игнорируем по причине проверки выше.
	dateIn, _ := time.Parse(date.FormatDate, taskData.Date)

	if now.After(dateIn) {
		if taskData.Repeat == "" {
			taskData.Date = fmt.Sprint(now.Format(date.FormatDate))
		} else {
			taskData.Date, err = date.NextDate(now, taskData.Date, taskData.Repeat)
			if err != nil {
				return taskData, err
			}
		}
	}

	return taskData, nil
}
