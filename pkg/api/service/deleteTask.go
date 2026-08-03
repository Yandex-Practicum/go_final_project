package service

import (
	"database/sql"
	"net/http"
)

func DeleteTask(id string) ([]byte, error, int) {
	// var res Task

	db, err := sql.Open("sqlite", PathDbManual())
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer db.Close()

	taska, err, status := GetTaskID(id)
	if err != nil {
		return nil, err, status
	}

	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", taska.ID))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	// jsByte, err := json.Marshal(res)
	// if err != nil {
	// 	return nil, err, http.StatusInternalServerError
	// }

	return []byte("{}"), nil, http.StatusOK

}
