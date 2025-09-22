package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"io"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	//Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, "Ошибка чтения тела запроса")
		return
	}

	//Десериализация
	if err := json.Unmarshal(body, &task); err != nil {
		writeError(w, "Ошибка парсинга JSON")
		return
	}

	//Валидация
	if task.Title == "" {
		writeError(w, "Заголовок не может быть пустым")
		return
	}

	//Установим дату "Сегодня"
	now := time.Now()
	layout := "20060102"
	nowDate := now.Format(layout)

	if task.Date == "" {
		task.Date = nowDate
	}

	//Проверка даты
	t, err := time.Parse(layout, task.Date)
	if err != nil {
		writeError(w, "Неверный формат даты")
		return
	}
	taskDate := t.Format(layout)

	//Обработка повторения задач
	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			writeError(w, "Неверное правило повторения")
			return
		}
		if taskDate < nowDate {
			nextDate, _ := NextDate(now, task.Date, task.Repeat)
			task.Date = nextDate
		}
	} else {
		if taskDate < nowDate {
			task.Date = nowDate
		}
	}

	//Сохраняем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "Ошибка при добавлении задачи")
		return
	}

	// Возвращаем ID новой задачи
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"id": id})
}
