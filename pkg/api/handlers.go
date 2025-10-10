package api

import (
	"fmt"
	"net/http"
	"time"
)

// NextDateHandler щбработчик для /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// NextDate остается публичным - не требует аутентификации
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Если now не указан, используем текущую дату
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			writeJSON(w, map[string]string{"error": fmt.Sprintf("Неверный формат now: %v", err)}, http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату
	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// TaskHandler основной обработчик для /api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
	}
}
