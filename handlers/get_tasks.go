package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/qavaleria/go_final_project/database"
	"net/http"
)

const Limit = 50

func HandleGetTasks(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		search := r.URL.Query().Get("search")
		tasks, err := database.GetTasks(db, search, Limit)
		if err != nil {
			http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
	}
}
