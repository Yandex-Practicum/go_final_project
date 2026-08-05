package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "modernc.org/sqlite"
)

func UpdateTask(w http.ResponseWriter, r *http.Request) ([]byte, error, int) {

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
		return nil, fmt.Errorf("Ошибка при проверке параметров"), http.StatusBadRequest
	}

	db, err := sql.Open("sqlite", PathDbManual())
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer db.Close()

	res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if count == 0 {
		return nil, fmt.Errorf("incorrect id for updating task"), http.StatusInternalServerError
	}

	return []byte("{}"), nil, http.StatusOK
}
