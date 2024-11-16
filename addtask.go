package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type response struct {
	Id    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func taskPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	t := task{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		http.Error(w, "Unable to decode JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if t.Title == "" {
		res := response{
			Error: "нет заголовка",
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	if t.Date == "" {
		t.Date = time.Now().Format(timeFormat)
	}

	dateTime, err := time.Parse(timeFormat, t.Date)
	if err != nil {
		res := response{
			Error: "неверный формат даты",
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	dateTime = dateTime.Truncate(24 * time.Hour)
	now := time.Now().Truncate(24 * time.Hour)
	if dateTime.Before(now) {
		if t.Repeat == "" {
			dateTime = time.Now()
			t.Date = dateTime.Format(timeFormat)
		} else {
			t.Date, err = NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				res := response{
					Error: "ошибка создания следующей даты",
				}
				json.NewEncoder(w).Encode(res)
				return
			}
		}
	}

	dbFile := os.Getenv("TODO_DBFILE")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		res := response{
			Error: "ошибка сервера",
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		res := response{
			Error: "ошибка сервера",
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		res := response{
			Error: "ошибка сервера",
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := response{
		Id: int(id),
	}
	json.NewEncoder(w).Encode(res)
}
