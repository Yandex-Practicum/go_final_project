package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

var response struct {
	Error string `json:"error,omitempty"`
}

// daysInMonth - функция, которая принимает год и месяц в качестве аргументов и возвращает количество дней в этом месяце.
func daysInMonth(year, month int) int {
	// Используем switch-case, чтобы определить количество дней в каждом месяце.
	switch month {
	// Месяцы с 31 дня.
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	// Месяцы с 30 дня.
	case 4, 6, 9, 11:
		return 30
	// Февраль.
	case 2:
		// Если год является высокосным, то в феврале 29 дней.
		if isLeapYear(year) {
			return 29
		}
		// Иначе в феврале 28 дней.
		return 28
	// Неверный номер месяца.
	default:
		return 0
	}
}

// isLeapYear - функция, которая принимает год в качестве аргумента и возвращает true, если год является высокосным, и false, если нет.
func isLeapYear(year int) bool {
	// Если год делится на 400 без остатка, то это высокосный год.
	if year%400 == 0 {
		return true
	}
	// Если год делится на 100 без остатка, то это не высокосный год.
	if year%100 == 0 {
		return false
	}
	// Если год делится на 4 без остатка, то это высокосный год.
	if year%4 == 0 {
		return true
	}
	// Иначе это не высокосный год.
	return false
}

func CreateDB(DBFile string) {
	db, err := sql.Open("sqlite3", DBFile)

	if err != nil {
		return
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			"id" INTEGER PRIMARY KEY,
			"date" CHAR(8) NOT NULL,
			"title" VARCHAR(255) NOT NULL,
			"comment" TEXT,
			"repeat" VARCHAR(128) NOT NULL
			)
		`)
	if err != nil {
		return
	}
	_, err = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_date ON sсheduler (date)
		`)
	if err != nil {
		return
	}
}

func main() {
	serverStart()
}
