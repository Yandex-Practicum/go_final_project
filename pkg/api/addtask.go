package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/VyacheZhuk/final-project/pkg/db"
)

type AnswerAddTask struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var ans AnswerAddTask

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		ans.Error = fmt.Sprintf("невозможно прочитать тело запроса %v", err.Error())
		writeJson(w, http.StatusBadRequest, ans)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		ans.Error = fmt.Sprintf("невозможно десериализовать тело запроса %v", err.Error())
		writeJson(w, http.StatusBadRequest, ans)
		return
	}

	if task.Title == "" {
		ans.Error = "не указан заголовок задачи"
		writeJson(w, http.StatusBadRequest, ans)
		return
	}

	err = checkDate(&task)
	if err != nil {
		ans.Error = err.Error()
		writeJson(w, http.StatusBadRequest, ans)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		ans.Error = err.Error()
		writeJson(w, http.StatusInternalServerError, ans)
		return
	}

	ans.ID = strconv.Itoa(int(id))
	writeJson(w, http.StatusCreated, ans)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("ошибка парсинга даты задачи (поле 'task.Date'): %v", err)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}
func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func afterNow(now, t time.Time) bool {
	return t.After(now)
}
