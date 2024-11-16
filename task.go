package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type pole struct {
	Tasks []task `json:"tasks"`
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
	}

	defer db.Close()
	getParam := r.FormValue("search")
	if getParam != "" {
		t, err := time.Parse("02.01.2006", getParam)
		var param string
		if err != nil {
			param = getParam
			result, err := db.Query("SELECT id,date,title,comment,repeat FROM scheduler WHERE comment LIKE :param OR title LIKE :param ORDER BY date ASC", sql.Named("param", "%"+param+"%"))

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
			}
			defer result.Close()

			tasks := []task{}
			for result.Next() {
				t := task{}
				result.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
				tasks = append(tasks, t)
			}
			defer result.Close()
			resp := pole{
				Tasks: tasks,
			}
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
			}
		} else {

			param = t.Format(timeFormat)
			result, err := db.Query("SELECT id,date,title,comment,repeat FROM scheduler WHERE date LIKE :param ORDER BY date ASC", sql.Named("param", "%"+param+"%"))

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
			}
			defer result.Close()

			tasks := []task{}
			for result.Next() {
				t := task{}
				result.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
				tasks = append(tasks, t)
			}
			defer result.Close()
			resp := pole{
				Tasks: tasks,
			}
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
			}
		}

	} else {
		result, err := db.Query("SELECT id,date,title,comment,repeat FROM scheduler ORDER BY date ASC")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
		}
		tasks := []task{}
		for result.Next() {
			t := task{}
			result.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
			tasks = append(tasks, t)
		}
		defer result.Close()
		resp := pole{
			Tasks: tasks,
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "ошибка сервера"})
		}
	}
}
