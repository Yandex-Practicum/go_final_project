package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/seregey-schr/go_final_project/pkg/utils"
)

//const dateFormat = "20060102"

// nextDateHandler обрабатывает GET-запросы /api/nextdate
// Пример запроса: /api/nextdate?now=20240126&date=20240126&repeat=y
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверка обязательных параметров
	if dstart == "" || repeat == "" {
		http.Error(w, "Не указаны параметры date или repeat", http.StatusBadRequest)
		return
	}

	// Парсим now, если не указан — используем текущую дату
	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
		// обнуляем время, оставляем только дату
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат параметра now, требуется "+dateFormat, http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату
	next, err := utils.NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка вычисления следующей даты: %v", err), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	fmt.Fprint(w, next)
}
