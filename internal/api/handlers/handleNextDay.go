package handlers

import (
	"fmt"
	"go-task-manager-final_project/internal/api"
	"go-task-manager-final_project/internal/services"
	"net/http"
	"time"
)

// nextDayHandler обрабатывает HTTP‑запрос на вычисление следующей даты по правилу повторения.
// Ожидает GET‑запрос с параметрами:
// - now (текущая дата в формате services.DateFormat);
// - date (стартовая дата в текстовом формате);
// - repeat (правило повторения, определяющее периодичность).
// Возвращает:
// - вычисленную дату в текстовом формате при успешном выполнении;
// - JSON с ошибкой при некорректных входных данных или сбое вычислений.
func handleNextDay(w http.ResponseWriter, r *http.Request) {

	// Получаем значения параметров из запроса
	nowString := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Парсим строку с текущей датой в тип time.Time
	// Используем формат, определённый в пакете services (services.DateFormat)
	now, err := time.Parse(services.DateFormat, nowString)
	if err != nil {
		// Если формат даты некорректен, возвращаем ошибку 400 Bad Request
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid 'now' date format",
		})
		return
	}

	// Вычисляем следующую дату с помощью функции из пакета services
	// Функция учитывает текущую дату, стартовую дату и правило повторения
	nextDate, err := services.NextDate(now, date, repeat)
	if err != nil {
		// При ошибке в вычислении даты возвращаем ошибку 400 с описанием
		api.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("failed to calculate next date: %v", err),
		})
		return
	}

	// Отправляем вычисленную дату в ответ в текстовом формате
	w.Write([]byte(nextDate))
}
