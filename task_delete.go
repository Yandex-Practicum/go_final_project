package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

)

func taskDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{Error: "ошибка удаления задачи - неправильный запрос со стороны клиента"})
		return
	}
	db, err := sql.Open("sqlite3", os.Getenv("TODO_DBFILE"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "ошибка удаления задачи - сервер не может устанвоить связь с базой данных"})
		return
	}
	defer db.Close()
	res, err := db.Exec("DELETE FROM scheduler WHERE id=:id", sql.Named("id", id))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "ошибка удаления задачи "})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "ошибка удаления задачи "})
		return
	}
	json.NewEncoder(w).Encode(response{})
}
