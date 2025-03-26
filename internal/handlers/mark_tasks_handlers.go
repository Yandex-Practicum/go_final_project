package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go_final_project-main/internal/model"
	"go_final_project-main/internal/nextdate"

	"github.com/jmoiron/sqlx"
)

func MarkTaskDoneHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		// Преобразование ID из строки в int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"Некорректный идентификатор"}`, http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли задача
		var task model.Task
		err = db.Get(&task, `SELECT * FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		// Если задача одноразовая (пустое поле repeat), удаляем ее
		if task.Repeat == "" {
			_, err = db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
			if err != nil {
				http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
				return
			}
			// Возвращаем пустой JSON
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}

		// Если задача периодическая, рассчитываем следующую дату
		nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		_, err = db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, nextDate, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
			return
		}

		// Возвращаем пустой JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func DeleteTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error":"Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		// Преобразование ID из строки в int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"Некорректный идентификатор"}`, http.StatusBadRequest)
			return
		}

		// Удаляем задачу из базы данных
		_, err = db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
		if err != nil {
			http.Error(w, `{"error":"Ошибка при удалении задачи"}`, http.StatusInternalServerError)
			return
		}

		// Возвращаем пустой JSON
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
