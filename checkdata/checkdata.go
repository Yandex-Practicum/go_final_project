package checkdata

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go_final_project/nextdate"
	"go_final_project/task"
)

func CheckData(req *http.Request) (task.Task, int, error) {
	var task task.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return task, 500, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return task, 500, err
	}
	// Поле title обязательно должно быть указано, иначе возвращаем ошибку.
	if task.Title == "" {
		return task, 400, errors.New(`{"error":"task title is not specified"}`)
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	// Если дата не указана, полю date присваивается сегодняшняя дата.
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	dateParse, err := time.Parse("20060102", task.Date)
	if err != nil {
		return task, 400, errors.New(`{"error":"incorrect date"}`)
	}
	var dateNew string
	if task.Repeat != "" {
		dateNew, err = nextdate.NextDate(now, task.Date, task.Repeat) // Проверяем корректность поля repeat
		if err != nil {
			return task, 400, err
		}
	}

	// Если поле date равен текущему дню, то date присваивается сегодняшний день.
	if task.Date == now.Format("20060102") {
		task.Date = now.Format("20060102")
	}

	// Если дата раньше сегодняшней, есть два варианта:
	// 1. Если поле repeat пусто, то полю date присваиваетя сегодняшняя дата.
	// 2. Иначе полю date присваиваетя следующая дата повторения, высчитанная ранее ф. NextDate.
	if dateParse.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = dateNew
		}
	}

	return task, 200, nil
}
