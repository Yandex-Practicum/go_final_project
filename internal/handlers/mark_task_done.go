package handlers

import (
	"database/sql"
	"fmt"
	"my_education/go/go_final_project/internal/logic"
	"net/http"
	"time"
)

// MarkTaskDoneHandler обрабатывает POST-запрос для отметки задачи как выполненной
func MarkTaskDoneHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем ID задачи из параметров запроса
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		// Загружаем задачу из базы данных
		var task Task
		query := "SELECT id, date, repeat FROM scheduler WHERE id = ?"
		err := db.QueryRow(query, taskID).Scan(&task.ID, &task.Date, &task.Repeat)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка при получении задачи: %v"}`, err), http.StatusInternalServerError)
			}
			return
		}

		// Если задача одноразовая (repeat пустой), удаляем её
		if task.Repeat == "" {
			deleteQuery := "DELETE FROM scheduler WHERE id = ?"
			_, err := db.Exec(deleteQuery, taskID)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"Ошибка при удалении задачи: %v"}`, err), http.StatusInternalServerError)
				return
			}

			// Возвращаем пустой JSON в случае успеха
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
			return
		}

		// Если задача повторяющаяся, вычисляем следующую дату
		now := time.Now()
		nextDate, err := logic.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		// Обновляем задачу в базе данных с новой датой
		updateQuery := "UPDATE scheduler SET date = ? WHERE id = ?"
		_, err = db.Exec(updateQuery, nextDate, taskID)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Ошибка при обновлении задачи: %v"}`, err), http.StatusInternalServerError)
			return
		}

		// Возвращаем пустой JSON в случае успеха
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}
