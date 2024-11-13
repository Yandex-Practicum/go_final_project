package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pwd/database"
	"pwd/internal/handler"
	"pwd/internal/nextdate"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, "недостаточно значений", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "недопустимый формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
}

// обрабатывает запросы для задач по методу запроса
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		PostTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

// хэндлер на запрос добавления задачи
func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task handler.Task // экзмпляр структуры со значениями

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "ошибка десериализации json", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Поле title обязательно", http.StatusBadRequest)
		return
	}

	// проверяем дату
	var date time.Time
	var err error
	if task.Date == "" {
		task.Date = time.Now().Format("20060102") // указзываем сегодняшнюю дату
	}
	date, err = time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, "некорректный формат даты", http.StatusBadRequest)
		return
	}

	now := time.Now()

	//  если дата меньше сегодняшней
	if date.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102") // Устанавливаем сегодняшнюю дату
		} else {
			// вычисляем следующую дату NextDate
			nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, "Ошибка вычисления следующей даты", http.StatusInternalServerError)
				return
			}
			task.Date = nextDate
		}
	}

	// добавляем задачу в базу данных
	id, err := database.AddTask(task)
	if err != nil {
		http.Error(w, "Ошибка при добавлении задачи в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)

}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

}
