package task

import (
	"encoding/json"
	"log"
	"net/http"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

func (h *Handler) handlePostTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO models.Task
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		utils.RespondWithError(w, utils.ErrInvalidJson)
		return
	}

	task, err := validateTask(&taskDTO)
	if err != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := h.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		utils.RespondWithError(w, utils.ErrDBInsert)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, utils.ErrGetTaskID)
		return
	}

	task.ID = id

	log.Printf("Задача добавлена: %+v\n", task)

	response := models.Response{ID: &task.ID}
	utils.SetJsonHeader(w)
	json.NewEncoder(w).Encode(response)
}
