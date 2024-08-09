package task_done

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

func (h *Handler) handlePostTaskDone(w http.ResponseWriter, r *http.Request) {
	stringId := r.URL.Query().Get("id")
	if len(stringId) == 0 {
		utils.RespondWithError(w, utils.ErrIDIsEmpty)
		return
	}
	id, err := strconv.ParseInt(stringId, 10, 64)
	if err != nil {
		utils.RespondWithError(w, utils.ErrIDIsEmpty)
		return
	}

	query := `SELECT 
    			id,
    			date,
    			repeat
			  FROM scheduler
			  WHERE id = ?`
	row := h.db.QueryRow(query, id)
	var task models.Task
	err = row.Scan(&task.ID, &task.Date, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, utils.ErrTaskNotFound)
			return
		}
		utils.RespondWithError(w, utils.ErrTaskParse)
		return
	}

	if len(task.Repeat) == 0 {
		deleteQuery := `DELETE FROM scheduler WHERE id = ?`
		_, deleteErr := h.db.Exec(deleteQuery, id)
		if deleteErr != nil {
			utils.RespondWithError(w, err.Error())
			return
		}
	} else {
		task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			utils.RespondWithError(w, err.Error())
			return
		}
		updateQuery := `UPDATE scheduler SET date = ? WHERE id = ?`
		_, err = h.db.Exec(updateQuery, task.Date, task.ID)
		if err != nil {
			utils.RespondWithError(w, err.Error())
			return
		}
	}

	utils.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
