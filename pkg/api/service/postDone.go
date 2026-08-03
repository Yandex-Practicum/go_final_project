package service

import (
	"database/sql"
	"net/http"
	"time"
)

func PostDone(id string, w http.ResponseWriter, r *http.Request) ([]byte, error, int) {
	var res string
	now := time.Now()

	taska, err, status := GetTaskID(id)
	if err != nil {
		return nil, err, status
	}

	if taska.Repeat == "" {
		jsByte, err, status := DeleteTask(id)
		if err != nil {
			return nil, err, status
		}
		return jsByte, nil, status
	} else {
		res, err = NextDate(now, taska.Date, taska.Repeat)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		db, err := sql.Open("sqlite", PathDbManual())
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		defer db.Close()

		_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("date", res),
			sql.Named("id", taska.ID))

		if err != nil {
			return nil, err, http.StatusInternalServerError
		}

		return []byte("{}"), nil, http.StatusOK

	}

}
