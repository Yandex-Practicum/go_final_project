package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ZhakiyaK/final_project/pkg/db"
)

// Добавление задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task
	var err error
	//Чтение тело запроса и декодирование
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {

		writeJsonError(w, fmt.Sprintf("incorrect data in JSON: %s", err), http.StatusBadRequest)

		return
	}

	//Проверка что поле task.Title не пустое
	if task.Title == "" {
		writeJsonError(w, "title required", http.StatusBadRequest)
		return
	}

	//Проверка на корректность полученное значение task.Date
	if err := checkDate(&task); err != nil {
		writeJsonError(w, fmt.Sprintf("Check date error: %s", err), http.StatusBadRequest)
		return
	}

	//Добавляем задачу в базу данных
	id, err := db.AddTask(&task)

	if err != nil {
		writeJsonError(w, fmt.Sprintf("DB add task error: + %s", err), http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]any{"id": id})
}

func checkDate(task *db.Task) error {

	now := time.Now().UTC().Truncate(24 * time.Hour)

	// Если task.Date пустая строка, то присваиваем ему текущее время
	if task.Date == "" {
		task.Date = now.UTC().Format("20060102")
	}
	// парсим дату
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	var next string

	// Если правило повторения указано, вычисляем следующую дату
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Проверка, что date больше чем now
	if afterNow(now, date) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}
