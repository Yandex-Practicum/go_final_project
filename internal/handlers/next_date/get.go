package next_date

import (
	"fmt"
	"net/http"
	"time"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, utils.ErrInvalidDateNowFormat, http.StatusBadRequest)
		return
	}

	nextDate, err := models.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, nextDate)
}
