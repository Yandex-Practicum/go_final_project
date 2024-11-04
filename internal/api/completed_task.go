package api

import (
	"database/sql"
	"net/http"
	"time"

	"go_final_project/internal/config"
	"go_final_project/internal/db"
	"go_final_project/internal/scheduler"
)

// CompleteTaskHandler обрабатывает POST-запросы к /api/task/done для выполнения задачи.
func CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")

	// Проверка, указан ли идентификатор
	if taskID == "" {
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// Подключение к базе данных
	dbConn, err := db.SetupDB()
	if err != nil {
		http.Error(w, `{"error": "Ошибка подключения к базе данных"}`, http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Получаем задачу из базы данных
	var task config.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = dbConn.Get(&task, query, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		// Если задача одноразовая, удаляем её из базы данных
		deleteQuery := `DELETE FROM scheduler WHERE id = ?`
		if _, err := dbConn.Exec(deleteQuery, taskID); err != nil {
			http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}
		// Возвращаем пустой JSON в случае успешного удаления
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		return
	}

	// Если задача периодическая, вычисляем следующую дату
	now := time.Now()
	nextDate, err := scheduler.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Обновляем дату задачи в базе данных
	updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err = dbConn.Exec(updateQuery, nextDate, taskID)
	if err != nil {
		http.Error(w, `{"error": "Ошибка при обновлении даты задачи"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON в случае успешного обновления
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}
