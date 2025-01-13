package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	nextDate, err := time.Parse("20060102", date)
	if err != nil {
		return "неверный формат даты", err
	}
	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat: %s", repeat)
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid days: %s", parts[1])
		}

		for {
			nextDate = nextDate.AddDate(0, 0, days)
			if now.Before(nextDate) {
				break
			}
		}
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid repeat: %s", repeat)
		}
		for {
			nextDate = nextDate.AddDate(1, 0, 0)
			if now.Before(nextDate) {
				break
			}
		}
	default:
		return "", fmt.Errorf("unknown repeat type: %s", parts[0])

	}

	return nextDate.Format("20060102"), nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Некорректная дата", http.StatusBadRequest)
		return
	}

	nextDates, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", nextDates)
}

func main() {

	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	println(appPath)

	dbFile := filepath.Join(appPath, "scheduler.db")
	_, err = os.Stat(dbFile)
	println(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install { //если файл БД не существует, то создаем его

		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY, date TEXT, title TEXT, comment TEXT, repeat TEXT(128))") // создаем таблицу с данными
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_date ON scheduler (date)")
		if err != nil {
			log.Fatal(err)
		}

	}

	webDir := "./web"
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}

}
