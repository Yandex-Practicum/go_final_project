package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/nextdate"
	"go_final_project/storage"
	"go_final_project/task"
)

func TaskDoneHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"missing id parameter"}`, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite", "./scheduler.db")
	if err != nil {
		http.Error(w, "error opening database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var taskID task.Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err = row.Scan(&taskID.Id, &taskID.Date, &taskID.Title, &taskID.Comment, &taskID.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"error scanning task: `+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	if taskID.Repeat == "" {
		ResponseStatus, err := storage.DeleteTask(db, id)
		if err != nil {
			http.Error(w, `{"error":"error deleting task: `+err.Error()+`"}`, ResponseStatus)
			return
		}
	} else {
		now := time.Now()
		nextDateStr, err := nextdate.NextDate(now, taskID.Date, taskID.Repeat)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		// Преобразование строки в time.Time и затем в нужный формат
		nextDate, err := time.Parse("20060102", nextDateStr)
		if err != nil {
			http.Error(w, `{"error":"error parsing next date: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		// Обновление даты задачи и создание нового запроса с телом
		taskID.Date = nextDate.Format("20060102")
		taskJson, err := json.Marshal(taskID)
		if err != nil {
			http.Error(w, `{"error":"error marshaling task: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		newReq, err := http.NewRequest(http.MethodPut, "", bytes.NewBuffer(taskJson))
		if err != nil {
			http.Error(w, `{"error":"error creating new request: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		newReq.Header.Set("Content-Type", "application/json")

		q := newReq.URL.Query()
		q.Add("id", id)
		newReq.URL.RawQuery = q.Encode()

		_, ResponseStatus, err := storage.UptadeTaskID(db, newReq)
		if err != nil {
			http.Error(w, `{"error":"error updating task: `+err.Error()+`"}`, ResponseStatus)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
