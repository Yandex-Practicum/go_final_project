package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

// GetNextDate returns the next date based on the provided parameters.
// It takes the current date and time, a specific date, and a repeat pattern as input.
// If an error occurs during the process, it sends an error response.
// Otherwise, it sends the next date as a plain text response.
func (tc *TaskController) GetNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nextDate, err := tc.uc.GetNextDate(now, date, repeat)
	if err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.PlainText(w, r, nextDate)
}
