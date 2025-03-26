package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/database"
	"github.com/qavaleria/go_final_project/models"
	"github.com/qavaleria/go_final_project/tasks"
	"log"
	"net/http"
)

func HandleUpdateTask(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Printf("Ошибка десериализации JSON: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if task.ID == "" {
			log.Printf("Не указан идентификатор задачи")
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		if err := tasks.ValidateTask(&task); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		err = database.UpdateTask(db, task)
		if err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
