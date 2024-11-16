package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func doTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "ошибка формирования отметки об выполнении задачи"})
		return
	}
	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "ошибка формирования отметки об выполнении задачи"})
		return
	}
	defer db.Close()
	res := db.QueryRow("SELECT repeat, date FROM scheduler WHERE id=:id", sql.Named("id", id))
	t := task{}
	res.Scan(&t.Repeat, &t.Date)
	if t.Date == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{Error: "не получилось получить данные о задаче"})
		return
	}
	if t.Repeat == "" {
		_, err = db.Exec("DELETE FROM scheduler WHERE id=:id", sql.Named("id", id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "ошибка формирования отметки об выполнении задачи"})
			return
		}
	} else {
		newDate, err := NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "ошибка формирования отметки об выполнении задачи"})
			return
		}
		_, err = db.Exec("UPDATE scheduler SET date=:date", sql.Named("date", newDate))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "ошибка формирования отметки об выполнении задачи"})
			return
		}
	}
	json.NewEncoder(w).Encode(response{})
}
