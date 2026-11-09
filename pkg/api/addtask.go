package api

import (
	"database/sql"
	"encoding/json"
	"github.com/Yandex-Practicum/go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func addTaskHandler(database *sql.DB, w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{
			"error": "Не указан заголовок задачи",
		}, http.StatusBadRequest)

		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(database, &task)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJson(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	}, http.StatusOK)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	today := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)

	// Если дата задачи уже прошла:
	if taskDate.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}
