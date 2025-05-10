package api

import (
	"encoding/json"
	"errors"
	"finalproject/pkg/db"
	"io"
	"net/http"
	"time"
)

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

func writeJSON(w http.ResponseWriter, response []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Handler_NextDate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	now := query.Get("now")
	date := query.Get("date")
	repeat := query.Get("repeat")

	if date == "" || repeat == "" {
		errHandler(w, "date and repeat parameters are required", errors.New(""))
		return
	}

	var nowTime time.Time
	var err error

	if now == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse("20060102", now)
		if err != nil {
			errHandler(w, "Invalid now parameter", err)
			return
		}
	}
	nextDate, err := nextDate(nowTime, date, repeat)
	if err != nil {
		errHandler(w, "", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
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

	writeJSON(w, res, http.StatusCreated)
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

	writeJSON(w, tasksByte, http.StatusOK)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		errHandler(w, "Forgot entered ID", errors.New(""))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		errHandler(w, "Task not found", err)
		return
	}

	taskByte, err := json.Marshal(task)
	if err != nil {
		errHandler(w, "Failed to encode JSON", err)
		return
	}

	writeJSON(w, taskByte, http.StatusOK)
}

func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if _, err := db.GetTask(task.ID); err != nil {
		errHandler(w, "No such Task", err)
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

	err = db.UpdateTask(&task)
	if err != nil {
		errHandler(w, "", err)
		return
	}

	emptyJSON, err := json.Marshal(map[string]string{"result": "ok"})
	if err != nil {
		errHandler(w, "", err)
	}

	writeJSON(w, emptyJSON, http.StatusOK)
}
