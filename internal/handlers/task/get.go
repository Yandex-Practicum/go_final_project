package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

type GetTaskResponseDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
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
    			title,
    			comment,
    			repeat
			  FROM scheduler
			  WHERE id = ?`
	row := h.db.QueryRow(query, id)
	var task models.Task
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, utils.ErrTaskNotFound)
			return
		}
		utils.RespondWithError(w, utils.ErrTaskParse)
		return
	}
	response := GetTaskResponseDTO{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	utils.SetJsonHeader(w)
	json.NewEncoder(w).Encode(response)
}
