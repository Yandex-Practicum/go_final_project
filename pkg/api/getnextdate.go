package api

import (
	"net/http"
	"time"

	"github.com/JKasus/go_final_project/pkg/constants"
	"github.com/JKasus/go_final_project/pkg/internal"
)

func getNextDayHandler(w http.ResponseWriter, r *http.Request) {
	var nowDate time.Time
	nowParam := r.URL.Query().Get("now")
	if nowParam != "" {
		var err error
		nowDate, err = time.Parse(constants.DateFormat, nowParam)
		if err != nil {
			nowDate = time.Now()
		}
	} else {
		nowDate = time.Now()
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := internal.NextDate(nowDate, date, repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
