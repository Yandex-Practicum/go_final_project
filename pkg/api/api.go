package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const DateFormat = "20060102"

// Init регистрирует все API обработчики
func Init() {
	InitAuth()

	http.HandleFunc("/api/task", AuthMiddleware(TaskHandler))
	http.HandleFunc("/api/tasks", AuthMiddleware(TasksHandler))
	http.HandleFunc("/api/task/done", AuthMiddleware(DoneTaskHandler))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/signin", SigninHandler)
}

// TaskHandler обрабатывает все запросы к /api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

// NextDateHandler обрабатывает запросы для /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	if date == "" {
		http.Error(w, "Параметр date обязателен", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "Параметр repeat обязателен", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

// writeJSONResponse отправляет JSON ответ с указанным статусом
func writeJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
	}
}

// writeJSONError отправляет JSON ответ с ошибкой
func writeJSONError(w http.ResponseWriter, message string, status int) {
	response := map[string]string{"error": message}
	writeJSONResponse(w, response, status)
}

// formatID преобразует int64 ID в строку
func formatID(id int64) string {
	return strconv.FormatInt(id, 10)
}
