package handler

import (
	"encoding/json"

	"net/http"
	"time"

	nextdate "github.com/askoren1/go_final_project/internal/next_date"
)

// Функция для обработки запроса на отметку задачи как выполненной
func (h *Handler) MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id") //Получение идентификатора задачи

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	//Получение информации о задаче
	task, err := h.repo.GetTaskByID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка получения задачи: " + err.Error()})
		return
	}

	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		// Одноразовая задача - удаляем
		err = h.repo.DeleteTask(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка удаления задачи: " + err.Error()})
			return
		}
	} else {
		// Периодическая задача - вычисляем следующую дату
		nowTime := time.Now()
		nextDate, err := nextdate.NextDate(nowTime, task.Date, task.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка вычисления следующей даты: " + err.Error()})
			return
		}

		err = h.repo.UpdateTask(idStr, nextDate, task.Title, task.Comment, task.Repeat)
		if err != nil {
			// Возвращаем JSON с ошибкой
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
	}

	// Возвращаем пустой JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}
