package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/internal/constants"
	"go_final_project/internal/error"
	"go_final_project/internal/task"
)

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {

	var t task.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	//Проверка поля title
	if t.Title == "" {
		error.JsonResponse(w, "Не указан заголовок задачи")
		return
	}

	now := time.Now()
	today := now.Format(constants.DateFormat)

	//Проверка поля date
	if t.Date == "" {
		t.Date = today
	} else {
		parseDate, err := time.Parse(constants.DateFormat, t.Date)
		if err != nil {
			error.JsonResponse(w, "Дата указана в неверном формате")
			return
		}

		parseDate = parseDate.Truncate(24 * time.Hour)
		now = now.Truncate(24 * time.Hour)

		//Проверка что дата меньше сегодняшнего числа и вычисление следующей даты повторения
		if parseDate.Before(now) {
			if t.Repeat == "" {
				t.Date = today
			} else {
				nextDate, err := NextDate(now, t.Date, t.Repeat)
				if err != nil {
					error.JsonResponse(w, "Неверный формат правила повторения")
					return
				}
				t.Date = nextDate
			}
		}
	}

	//Добавляем задачу в db
	id, err := h.repo.AddTask(t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Формируем ответ
	response := map[string]string{"id": fmt.Sprintf("%d", id)}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
