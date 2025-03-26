package handlers

import (
	"database/sql"
	"encoding/json"
	"go_final_project/utils"
	"net/http"
	"time"
)

func NextDateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := r.FormValue("now")
		date := r.FormValue("date")
		repeat := r.FormValue("repeat")

		nowTime, err := time.Parse("20060102", now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nextDate, err := utils.NextDate(nowTime, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{"nextDate": nextDate})
	}
}
