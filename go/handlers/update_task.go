package handlers

import (
	"encoding/json"
	"final-project-bronner/go/models"
	"final-project-bronner/go/nextdate"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t models.Task
	now := time.Now()

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t.Id == "" {
		responseError(w, "id не может быть пустым", http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(t.Id); err != nil {
		responseError(w, "id должен быть числом", http.StatusBadRequest)
		return
	}
	if t.Title == "" {
		responseError(w, "title не может быть пустым", http.StatusBadRequest)
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

	err = h.db.UpdateTask(t)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(struct{}{})
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
