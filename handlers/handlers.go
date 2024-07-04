package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"main.go/api"
	"main.go/database"
	"main.go/model"
)

const layoutDate string = "20060102" // формат даты

func HandleNextDate(w http.ResponseWriter, r *http.Request) {

	nowRaw := r.URL.Query().Get("now")
	now, err := time.Parse(layoutDate, nowRaw)
	if err != nil {
		responseWithError(w, err)
	}
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	result, err := api.NextDate(now, date, repeat, layoutDate)
	out := result
	if err != nil {
		out = err.Error()
	}

	w.Write([]byte(out))
}

func responseWithError(w http.ResponseWriter, err error) {
	fmt.Printf("%v\n", err)
	json.NewEncoder(w).Encode(model.ResponseError{Error: err.Error()})

}

func HandleTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		switch r.Method {

		case http.MethodPost:
			var task model.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				responseWithError(w, errors.New("Ошибка десериализации JSON"))
				return
			}

			if task.Title == "" {
				responseWithError(w, errors.New("Не указан заголовок задачи"))
				return
			}
			task.Date, err = api.GetNextDate(task, layoutDate)
			if err != nil {
				responseWithError(w, errors.New("Ошибка получения следующей даты"))
				return
			}

			id, err := database.AddTask(db, task)
			if err != nil {
				responseWithError(w, errors.New("Ошибка добавления задачи"))
				return
			}

			response := map[string]interface{}{
				"id": id,
			}
			json.NewEncoder(w).Encode(response)

		case http.MethodGet:

			tasks, err := database.GetTasks(db)
			if err != nil {
				responseWithError(w, errors.New("Ошибка получения задач"))
			}

			result := model.Tasks{Tasks: tasks}

			if tasks == nil {
				result = model.Tasks{Tasks: []model.Task{}}

			}
			json.NewEncoder(w).Encode(result)

		case http.MethodPut:

			var task model.Task

			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				responseWithError(w, err)
				return
			}

			if task.ID == "" {
				responseWithError(w, errors.New("Не указан ID задачи"))
				return
			}

			if task.Title == "" {
				responseWithError(w, errors.New("Не указан заголовок задачи"))
				return
			}

			task.Date, err = api.GetNextDate(task, layoutDate)
			if err != nil {
				responseWithError(w, errors.New("Ошибка получения следующей даты"))
				return
			}

			if err = database.UpdateTask(db, task); err != nil {
				responseWithError(w, err)
				return

			}

			json.NewEncoder(w).Encode(map[string]interface{}{})

		case http.MethodDelete:

			id := r.URL.Query().Get("id")
			if id == "" {
				responseWithError(w, errors.New("Задача не найдена"))
				return
			}

			if err := database.DeleteTask(db, id); err != nil {
				responseWithError(w, err)
				return

			}

			json.NewEncoder(w).Encode(map[string]interface{}{})

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func HandleTaskID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			responseWithError(w, errors.New("Задача не найдена"))
			return
		}
		task, err := database.GetTask(db, id)
		if err != nil {
			responseWithError(w, errors.New("Ошибка получения задачи"))
			return
		}
		json.NewEncoder(w).Encode(task)
	}
}

func HandleTaskDone(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var task model.Task

		id := r.URL.Query().Get("id")

		if id == "" {
			responseWithError(w, errors.New("Задача не найдена"))
			return
		}
		task, err := database.GetTask(db, id)
		if err != nil {
			responseWithError(w, errors.New("Ошибка получения задачи"))
			return
		}

		if task.Repeat != "" {

			nextDate, err := api.NextDate(time.Now(), task.Date, task.Repeat, layoutDate)
			if err != nil {
				responseWithError(w, err)
				return
			}
			task.Date = nextDate

			query := `UPDATE scheduler SET date = ? WHERE id = ?`

			res, err := db.Exec(query, task.Date, task.ID)
			if err != nil {
				responseWithError(w, err)
				return
			}

			rows, err := res.RowsAffected()
			if err != nil {
				responseWithError(w, err)
				return
			}

			if rows != 1 {
				responseWithError(w, errors.New("expected to affect 1 row"))
				return
			}
		} else {
			if err := database.DeleteTask(db, id); err != nil {
				responseWithError(w, err)
				return

			}

		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
