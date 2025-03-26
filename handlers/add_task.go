package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/database"
	"github.com/qavaleria/go_final_project/models"
	"log"
	"net/http"
)

// HandleAddTask обработчик для добавления задачи
func HandleAddTask(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Printf("Ошибка десериализации JSON: %v", err)
			http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		id, err := database.AddTask(db, task)
		if err != nil {
			log.Printf("Ошибка добавления задачи: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}
