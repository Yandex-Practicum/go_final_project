package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/database"
	"github.com/qavaleria/go_final_project/tasks"
	"log"
	"net/http"
	"time"
)

// HandleMarkTaskDone обработчик для пометки задачи выполненной
func HandleMarkTaskDone(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		task, err := database.GetTaskByID(db, id)
		if err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			}
			return
		}

		if task.Repeat == "" {
			// Удаляем одноразовую задачу
			err := database.DeleteTask(db, id)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		} else {
			// Рассчитываем следующую дату для периодической задачи
			now := time.Now()
			nextDate, err := tasks.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Printf("Ошибка вычисления следующей даты: %v", err)
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}

			// Обновляем задачу с новой датой
			task.Date = nextDate
			err = database.UpdateTask(db, task)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
