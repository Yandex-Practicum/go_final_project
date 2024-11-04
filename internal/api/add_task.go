package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/internal/config"
	"go_final_project/internal/db"
	"go_final_project/internal/scheduler"
)

// TaskHandler обрабатывает запросы к /api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	default:
		http.Error(w, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
	}
}

// AddTaskHandler добавляет новую задачу в базу данных
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task config.Task

	// Десериализация JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error": "ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Обработка даты
	now := time.Now().Truncate(24 * time.Hour)

	// Устанавливаем сегодняшнюю дату, если не указана
	if task.Date == "" || task.Date == "today" {
		task.Date = now.Format(config.TimeFormat)
	}

	dateTask, err := time.Parse(config.TimeFormat, task.Date)
	if err != nil {
		http.Error(w, `{"error": "недопустимый формат даты"}`, http.StatusBadRequest)
		return
	}

	// Проверяем дату задачи
	if dateTask.Before(now) {
		if task.Repeat == "" {
			// Если дата меньше сегодняшней и нет правила повторения, устанавливаем текущую дату
			task.Date = now.Format(config.TimeFormat)
		} else {
			// Если указано правило повторения, вычисляем следующую дату
			nextDate, err := scheduler.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	}

	// Добавление задачи в базу данных
	dbConn, err := db.SetupDB()
	if err != nil {
		http.Error(w, `{"error": "ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := dbConn.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error": "ошибка добавления задачи в базу данных"}`, http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error": "ошибка получения идентификатора добавленной задачи"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем идентификатор задачи в формате JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}
