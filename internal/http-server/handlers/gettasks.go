package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	sl "github.com/wissio/go_final_project/internal/lib/logger/slog"
	"github.com/wissio/go_final_project/internal/models"
	"github.com/wissio/go_final_project/internal/storage/sqlite"
)

func GetTasks(log *slog.Logger, s *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTasks"
		log = log.With(slog.String("op", op))

		searchParam := r.URL.Query().Get("search")
		dateParam := ""
		limitParam := r.URL.Query().Get("limit")
		var limit int

		if limitParam != "" {
			var err error
			limit, err = strconv.Atoi(limitParam)
			if err != nil || limit < 0 {
				limit = 0
			}
		}

		if isDotLayout(searchParam) {
			dateParam = changeLayout(searchParam)
			searchParam = ""
		}

		tasks, err := s.GetTasks(dateParam, searchParam, limit)
		if err != nil {
			http.Error(w, `{"error":"Error retrieving tasks"}`, http.StatusInternalServerError)
			log.Error("Error getting tasks", sl.Err(err))
			return
		}

		response := map[string]interface{}{
			"tasks": tasks,
		}

		if tasks == nil {
			response["tasks"] = []models.Task{}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error("Error writing response", sl.Err(err))
		}
	}
}

func isDotLayout(str string) bool {
	_, err := time.Parse(DateDotLayout, str)
	return err == nil
}

func changeLayout(str string) string {
	date, _ := time.Parse(DateDotLayout, str)
	return date.Format(DateLayout)
}
