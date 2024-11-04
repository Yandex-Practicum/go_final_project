package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/internal/config"
	"go_final_project/internal/db"
	"go_final_project/internal/scheduler"
)

// GetTaskHandler обрабатывает GET-запрос для получения задачи по идентификатору.
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	// Проверка, указан ли идентификатор
	if taskID == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Подключение к базе данных
	dbConn, err := db.SetupDB()
	if err != nil {
		http.Error(w, `{"error": "ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Получение задачи по идентификатору
	var task config.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = dbConn.Get(&task, query, taskID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем задачу в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTaskHandler обрабатывает PUT-запрос для обновления задачи по идентификатору.
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task config.Task

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error": "Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if task.Id == "" {
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Подключение к базе данных
	dbConn, err := db.SetupDB()
	if err != nil {
		http.Error(w, `{"error": "Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Получаем текущую задачу из базы данных
	var currentTask config.Task
	query := `SELECT * FROM scheduler WHERE id=?`
	err = dbConn.Get(&currentTask, query, task.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
		return
	}

	// Обновляем поля только если они не пустые
	if task.Title != "" {
		currentTask.Title = task.Title
	}
	if task.Comment != "" {
		currentTask.Comment = task.Comment
	}
	if task.Repeat != "" {
		currentTask.Repeat = task.Repeat
	}

	// Обработка даты
	now := time.Now()
	nowStr := now.Format(config.TimeFormat) // Строковое представление текущей даты

	// Устанавливаем сегодняшнюю дату, если не указана
	if task.Date == "" || task.Date == "today" {
		task.Date = nowStr
	} else {
		dateTask, err := time.Parse(config.TimeFormat, task.Date)
		if err != nil {
			http.Error(w, `{"error": "недопустимый формат даты"}`, http.StatusBadRequest)
			return
		}

		// Проверяем дату задачи
		if dateTask.Before(now) {
			if task.Repeat == "" {
				// Если дата меньше сегодняшней и нет правила повторения, устанавливаем текущую дату
				task.Date = nowStr
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
	}

	// Обновление задачи в базе данных
	updateQuery := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := dbConn.Exec(updateQuery, task.Date, currentTask.Title, currentTask.Comment, currentTask.Repeat, task.Id)
	if err != nil {
		http.Error(w, `{"error": "Ошибка обновления задачи"}`, http.StatusInternalServerError)
		return
	}

	// Проверка, было ли обновлено хотя бы одно поле
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Успешное обновление
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}
