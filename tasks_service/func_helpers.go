package tasks_service

import (
	"database/sql"
	"net/http"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const dateFormat = "20060102"

func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func containsInt(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func findNextMonth(curMonth int, months []int) int {
	for _, month := range months {
		if month >= curMonth {
			return month
		}
	}
	return months[0]
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}

func compareDates(t1, t2 time.Time) (time.Time, time.Time) {
	date1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	date2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	return date1, date2
}
