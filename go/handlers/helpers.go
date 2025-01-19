package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func responseError(w http.ResponseWriter, error string, code int) {
	escapedError, err := json.Marshal(error)
	if err != nil {
		http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{\"error\": %s}", string(escapedError))
}

func dateIsNotNow(now, date time.Time) bool {
	year, month, day := date.Date()
	yearNow, monthNow, dayNow := now.Date()
	return dayNow != day || monthNow != month || yearNow != year
}
