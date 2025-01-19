package handlers

import (
	"encoding/json"
	"final-project-bronner/go/models"
	"final-project-bronner/go/nextdate"
	"net/http"
	"time"
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t models.Task
	now := time.Now()
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		responseError(w, "Title не может быть пустым", http.StatusBadRequest)
		return
	}

	if t.Date == "" {
		t.Date = now.Format(DateFormat)
	}
	dateParse, err := time.Parse(DateFormat, t.Date)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !dateParse.After(now) && dateIsNotNow(now, dateParse) {
		var assignDateBuf string
		if t.Repeat == "" {
			assignDateBuf = now.Format(DateFormat)
		} else {
			nextDate, err := nextdate.NextDate(now, t.Date, t.Repeat)
			if err != nil {
				responseError(w, err.Error(), http.StatusBadRequest)
				return
			}
			assignDateBuf = nextDate
		}
		t.Date = assignDateBuf
	}

	resultId, err := h.db.AddTask(t)

	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := struct {
		Id int `json:"id"`
	}{Id: resultId}
	result, err := json.Marshal(res)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
