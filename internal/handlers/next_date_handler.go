package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go_final_project/internal/utils"
)

type NextDateHandler struct {
}

func NewNextDateHandler() *NextDateHandler {
	return &NextDateHandler{}
}

func (h *NextDateHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleGet(w, r)
		default:
			http.Error(w, utils.ErrUnsupportedMethod, http.StatusMethodNotAllowed)
		}
	}
}

func (h *NextDateHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(utils.ParseDateFormat, nowStr)
	if err != nil {
		http.Error(w, utils.ErrInvalidDateNowFormat, http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, nextDate)
}
