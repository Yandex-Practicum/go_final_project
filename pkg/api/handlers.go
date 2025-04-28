package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"final_sprint/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализация
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		writeError(w, "Не указан заголовок задачи")
		return
	}

	// Проверка и корректировка даты
	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	// Добавляем задачу
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	// Успешный ответ
	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	now := time.Now().UTC()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		var timeErr *time.ParseError
		if errors.As(err, &timeErr) {
			if timeErr.Message == "day out of range" {
				return fmt.Errorf("несуществующий день в дате")
			}
		}
		return fmt.Errorf("дата содержит недопустимое значение: %q не существует", task.Date)
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("неверное правило повторения: %v", err)
		}

		if afterNow(now, t) {
			task.Date = next
		}

		return nil
	}

	if afterNow(now, t) {
		task.Date = now.Format(DateFormat)
	}

	return nil
}

// Сравнение
func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}

// Отправка успешного JSON-ответа
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// Отправка ошибки в формате JSON
func writeError(w http.ResponseWriter, errMsg string) {
	writeJSON(w, map[string]string{"error": errMsg})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "НЕ JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", task.ID)})
}

// Обработчик GET /api/task
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан id"})
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		writeJSON(w, map[string]string{"error": "Некорректный формат идентификатора"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, task)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		// Одноразовая задача — удаляем
		err := db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]interface{}{})
		return
	}

	// Периодическая задача — обновляем дату
	now, _ := ParseNowOrDefault("")
	newDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	err = db.UpdateDate(newDate, id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан ID задачи"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, struct{}{})
}
