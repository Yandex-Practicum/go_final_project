package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"strconv"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(w http.ResponseWriter, r *http.Request) ([]byte, error, int) {
	var id int64
	var task Task

	byteBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	err = json.Unmarshal(byteBody, &task)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if task.Title == "" {
		return nil, fmt.Errorf("Не указан заголовок задачи"), http.StatusBadRequest
	}
	err = CheckDate(&task)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при проверке параметров."), http.StatusBadRequest
	}

	db, err := sql.Open("sqlite", PathDbManual())
	if err != nil {
		log.Println(err)
		return nil, err, http.StatusInternalServerError
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}

	type response struct {
		ID string `json:"id"`
	}
	var resp response

	resp.ID = strconv.Itoa(int(id))

	jsByte, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Ошибка подготовки ответа", http.StatusInternalServerError)
		return nil, err, http.StatusInternalServerError
	}

	return jsByte, nil, http.StatusCreated

}
