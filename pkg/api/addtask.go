package api

import (
	"Final/pkg/db"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/http"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка сохранения задачи: " + err.Error()})
		return
	}

	writeJson(w, map[string]interface{}{"id": id})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("некорректная дата: %v", err)
	}

	startOfDay := func(t time.Time) time.Time {
		y, m, d := t.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %v", err)
		}

		if !t.After(startOfDay(now)) {
			task.Date = next
		}
	} else {
		if t.Before(startOfDay(now)) {
			task.Date = now.Format("20060102")
		}
	}

	return nil
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		writeJson(w, map[string]string{"error": "некорректный JSON"})
		return
	}
	if t.ID == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	if t.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}
	err = checkRepeat(t.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	err = checkDate(&t)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&t)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJson(w, struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJson(w, struct{}{})
}

func taskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "Ошибка удаления"})
			return
		}
		writeJson(w, struct{}{})
		return
	}
	nextdata, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось раcчитать следующую дату"})
		return
	}
	err = db.UpdateDate(nextdata, id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось обновить дату"})
		return
	}
	writeJson(w, struct{}{})
}

func checkRepeat(repeat string) error {
	if repeat == "" {
		return nil
	}

	parts := strings.Fields(repeat)
	if len(parts) != 2 {
		return errors.New("некорректный repeat")
	}

	unit := parts[0]
	if unit != "d" && unit != "w" && unit != "m" && unit != "y" {
		return errors.New("некорректная единица repeat")
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil || num <= 0 {
		return errors.New("некорректное число repeat")
	}

	return nil
}
