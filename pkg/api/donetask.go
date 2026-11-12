package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
	"go_final_project/pkg/nextdate"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{
			"error": "Метод не поддерживается",
		})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "Задача не найдена",
		})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error(),
			})
			return
		}
	} else {
		now := time.Now()
		next, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error(),
			})
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}

	writeJSON(w, map[string]string{})
}
