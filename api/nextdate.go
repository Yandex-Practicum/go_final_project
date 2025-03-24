package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// AddTaskHandler обрабатывает POST-запросы для добавления задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
		return
	}

	var task Task
	// Декодируем JSON из запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if task.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверяем поле repeat
	if task.Repeat != "" && !isValidRepeat(task.Repeat) {
		http.Error(w, `{"error": "Недопустимое значение для поля repeat"}`, http.StatusBadRequest)
		return
	}

	// Проверяем формат даты
	now := time.Now().Truncate(24 * time.Hour) // Убираем время, оставляем только дату
	var taskDate time.Time
	var err error

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	taskDate, err = time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error": "Дата представлена в неправильном формате"}`, http.StatusBadRequest)
		return
	}

	// Проверяем, не меньше ли дата сегодняшнего числа
	if taskDate.Before(now) {
		taskDate = now // Устанавливаем на сегодняшнюю дату
	} else if taskDate.Equal(now) {
		// Если дата равна текущей, оставляем её
		taskDate = now
	}

	// Если правило повторения указано как еженедельное или ежемесячное, возвращаем ошибку
	if task.Repeat != "" && (strings.HasPrefix(task.Repeat, "w") || strings.HasPrefix(task.Repeat, "m")) {
		http.Error(w, `{"error": "Неподдерживаемое правило повторения"}`, http.StatusBadRequest)
		return
	}

	// Устанавливаем task.Date на основе taskDate
	task.Date = taskDate.Format("20060102")

	// Добавляем задачу в базу данных
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error": "Ошибка при добавлении задачи"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error": "Ошибка при получении идентификатора задачи"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Ошибка при формировании ответа"}`, http.StatusInternalServerError)
		return
	}
}

// NextDate вычисляет следующую дату на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("некорректная дата")
	}

	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	switch {
	case repeat == "y":
		// Ежегодное повторение
		nextDate := startDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil

	case strings.HasPrefix(repeat, "d "):
		// Проверка формата d <число>
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("недопустимое количество дней")
		}

		// Найти следующую дату с учетом повторения
		nextDate := startDate.AddDate(0, 0, days)

		// Если следующая дата равна или совпадает с текущей, вернуть ее
		if nextDate.Equal(now) || nextDate.After(now) {
			return nextDate.Format("20060102"), nil
		}

		// Продолжаем добавлять дни, пока не найдем следующую дату
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

		// Если следующая дата равна или позже текущей, вернуть ее
		if nextDate.After(now) {
			return nextDate.Format("20060102"), nil
		}

		return "", errors.New("нет следующей даты")
	default:
		return "", errors.New("неподдерживаемый формат")
	}
}

// NextDateHandler обрабатывает запросы для получения следующей даты
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "некорректная дата now", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, nextDate)
}
