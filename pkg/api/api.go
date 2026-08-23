package api

import (
	"net/http"
	"time"
)

// Init регистрирует HTTP-обработчики API.
func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/signin", signInHandler)
}


// nextDayHandler обрабатывает запрос на вычисление следующей даты выполнения задачи.
//
// Если параметр now не указан, используется текущая дата.
// При успешном выполнении возвращается следующая дата в формате 20060102.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var nowTime time.Time
	var err error

	// если дата now не указана, используем текущую дату.
	if now == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse(dateLayout, now)
		if err != nil {
			http.Error(w, "invalid date format, expected YYYYMMDD", http.StatusBadRequest)
			return
		}
	}

	if date == "" {
		http.Error(w, "date is required", http.StatusBadRequest)
		return
	}
	
	if repeat == "" {
		http.Error(w, "repeat is required", http.StatusBadRequest)
		return
	}

	result, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)
	
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}