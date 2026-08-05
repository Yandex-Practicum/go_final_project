package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetTaskID(id string) (Task, error, int) {

	if id == "" {
		return Task{}, fmt.Errorf("Поле id не может быть пустым."), http.StatusBadRequest
	}

	db, err := sql.Open("sqlite", PathDbManual())
	if err != nil {
		return Task{}, err, http.StatusInternalServerError
	}
	defer db.Close()
	var res Task

	intStr, err := strconv.Atoi(id)

	err = db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", intStr)).Scan(&res.ID, &res.Date, &res.Title, &res.Comment, &res.Repeat)

	if err != nil {
		log.Println(err)
		return Task{}, err, http.StatusNotFound
	}

	return res, nil, http.StatusOK

}
