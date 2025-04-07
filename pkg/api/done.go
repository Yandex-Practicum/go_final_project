package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Stanislav-id05/go_final_project/pkg/db"
)

func TaskDone(w http.ResponseWriter, r *http.Request) {

	idGet := r.URL.Query().Get("id")
	if idGet == "" {
		http.Error(w, `{"error":"Нет индентификатора"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idGet)
	if err != nil {
		http.Error(w, `{"error":"Неверный идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(strconv.Itoa(id))
	if err != nil {
		http.Error(w, `{"error":"Результат поиска задачи отсутствует"}`, http.StatusNotFound)
		return
	}

	if task.Repeat != "" {
		if currentDate, err := time.Parse("20060102", task.Date); err == nil {
			if nextDate, err := NextDate(currentDate, task.Date, task.Repeat); err == nil {
				task.Date = nextDate
			} else {
				http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, `{"error":"Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
	} else if err = db.DeleteTask(idGet); err != nil {
		http.Error(w, `{"error":"Удаление невозможно"}`, http.StatusInternalServerError)
		return
	} else {
		w.Write([]byte("{}"))
		return
	}

	if err = db.UpdateTask(task); err != nil {
		http.Error(w, `{"error":""}`, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{}"))
}
