package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	utils "github.com/falsefood/go_final_project/internal"
)

func CreateTaskHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task utils.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if task.Title == "" {
		utils.WriteError(w, http.StatusBadRequest, "Заголовок нужно обязательно указать")
		return
	}

	today := time.Now()
	todayStr := today.Format("20060102")

	if task.Date == "" || task.Date == "today" || task.Date == todayStr {
		task.Date = todayStr
	} else {
		parsedDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Неверный формат даты")
			return
		}
		if parsedDate.Before(today) {
			if task.Repeat == "" {
				task.Date = todayStr
			} else {
				nextDate, err := nextDate(today, task.Date, task.Repeat)
				if err != nil {
					utils.WriteError(w, http.StatusBadRequest, err.Error())
					return
				}
				task.Date = nextDate
			}
		}
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при создании задачи")
		return
	}

	taskID, _ := result.LastInsertId()
	response := utils.Response{ID: taskID}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Ошибка при отправке данных")
		return
	}
}
