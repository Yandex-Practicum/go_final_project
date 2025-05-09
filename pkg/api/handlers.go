package api

import (
	"encoding/json"
	"errors"
	"finalproject/pkg/db"
	"io"
	"net/http"
	"time"
)

func Handler_NextDate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	now := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	if date == "" || repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date and repeat parameters are required"))
		return
	}

	var nowTime time.Time
	var err error

	if now == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse("20060102", now)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid now parameter " + err.Error()))
			return
		}
	}
	nextDate, err := nextDate(nowTime, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

func errHandler(w http.ResponseWriter, message string, err error) {
	errStruct := struct {
		Error string `json:"error"`
	}{
		Error: "message: " + message + "; error: " + err.Error(),
	}
	errBody, _ := json.Marshal(errStruct)

	w.WriteHeader(http.StatusBadRequest)
	w.Write(errBody)
}

func AddTaskHandle(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		errHandler(w, "Could not read body request", err)
		return
	}

	var task db.Task
	if err = json.Unmarshal(bodyByte, &task); err != nil {
		errHandler(w, "Error when decoding body", err)
		return
	}

	if task.Title == "" {
		errHandler(w, "Empty title field", errors.New("title is required"))
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		if _, err := time.Parse("20060102", task.Date); err != nil {
			errHandler(w, "Incorrect date format (expected YYYYMMDD)", err)
			return
		}
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		errHandler(w, "Incorrect date", err)
		return
	}

	if afterNow(now, t) {
		if task.Repeat == "" || task.Repeat == "d 1" {
			task.Date = now.Format("20060102")
		} else {
			next, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				errHandler(w, "", err)
				return
			}

			task.Date = next
		}
	}

	id, err := db.AddTask(task)
	if err != nil {
		errHandler(w, "", err)
		return
	}

	successResponse := struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}

	res, err := json.Marshal(successResponse)
	if err != nil {
		errHandler(w, "Failed to encode JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tip string

	search := r.URL.Query().Get("search")
	if t, err := time.Parse("02.01.2006", search); err == nil {
		search = t.Format("20060102")
		tip = "time"
	} else {
		tip = "default"
	}

	tasks, err := db.Tasks(50, search, tip)
	if err != nil {
		errHandler(w, "", err)
		return
	}

	var resp db.TaskResp
	if len(tasks) == 0 {
		resp = db.TaskResp{
			Tasks: []*db.Task{},
		}
	} else {
		resp = db.TaskResp{
			Tasks: tasks,
		}
	}

	tasksByte, err := json.Marshal(resp)
	if err != nil {
		errHandler(w, "Failed to encode JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tasksByte)
}
