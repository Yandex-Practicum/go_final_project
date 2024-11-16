package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"net/http"
	"os"
)

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")
	fmt.Println(id)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{Error: "ошибка запроса, не указан id в параметрах запроса"})
		return
	}
	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
		return
	}
	defer db.Close()
	res := db.QueryRow("SELECT * FROM scheduler WHERE id=?", id)

	t := task{}

	res.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if t.Date == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{Error: "ошибка сервера, задачи с таким id не существует"})
		return
	}

	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
		return
	}
}
func taskPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	t := task{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if t.ID == "" || t.Date == "" || t.Title == "" || t.Comment == "" || t.Repeat == "" {
		json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
		return
	}

	r.Body.Close()
	if err != nil {
		json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
		return
	}
	if t.ID == "" {
		json.NewEncoder(w).Encode(response{Error: "нет данных об id задачи"})
		return
	}
	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
		return
	}
	defer db.Close()
	res, err := db.Exec("UPDATE scheduler SET date=:date, title=:title, comment=:comment,repeat=:repeat WHERE id=:id", sql.Named("title", t.Title), sql.Named("date", t.Date), sql.Named("repeat", t.Repeat), sql.Named("comment", t.Comment), sql.Named("id", t.ID))
	rows, _ := res.RowsAffected()
	if rows != 1 {
		json.NewEncoder(w).Encode(response{Error: "ошибка обновления задачи"})
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(response{Error: "ошибка обновления данных о задаче"})
		return
	}
	json.NewEncoder(w).Encode(response{})
}
