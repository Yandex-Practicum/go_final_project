package handlers

import (
	"encoding/json"
	"final-project-bronner/go/models"
	"final-project-bronner/go/nextdate"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) MakeTaskHasDone(w http.ResponseWriter, r *http.Request) {
	var id int
	now := time.Now()
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := h.db.GetTask(id)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Repeat == "" {
		err := h.db.DeleteTask(id)
		if err != nil {
			responseError(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := json.Marshal(struct{}{})
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return
	}

	nextDate, err := nextdate.NextDate(now, t.Date, t.Repeat)
	if err != nil {
		err := h.db.DeleteTask(id)
		if err != nil {
			responseError(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := json.Marshal(models.Task{})
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
		return
	}
	t.Date = nextDate
	err = h.db.UpdateTask(t)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
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
