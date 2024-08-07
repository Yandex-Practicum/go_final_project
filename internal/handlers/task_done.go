package handlers

import (
	"database/sql"
	"errors"
	"go_final_project/internal/models"
	"net/http"
	"strconv"
	"time"
)

func TaskDoneHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePostTaskDone(w, r, db)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}
}

func handlePostTaskDone(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	stringId := r.URL.Query().Get("id")
	if len(stringId) == 0 {
		respondWithError(w, "Не указан идентификатор")
		return
	}
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		respondWithError(w, "Не указан идентификатор")
		return
	}

	query := `SELECT 
    			id,
    			date,
    			repeat
			  FROM scheduler
			  WHERE id = ?`
	row := db.QueryRow(query, id)
	var task models.Task
	err = row.Scan(&task.ID, &task.Date, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, "Задача не найдена")
			return
		}
		respondWithError(w, "Ошибка разбора задач из базы данных")
		return
	}

	if len(task.Repeat) == 0 {
		deleteQuery := `DELETE FROM scheduler WHERE id = ?`
		_, deleteErr := db.Exec(deleteQuery, id)
		if deleteErr != nil {
			respondWithError(w, err.Error())
			return
		}
	} else {
		task.Date, err = models.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}
		updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
		_, err = db.Exec(updateQuery, task.Date, task.ID)
		if err != nil {
			respondWithError(w, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
