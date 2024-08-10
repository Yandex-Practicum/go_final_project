package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_final_project/internal/models"
	"go_final_project/internal/utils"
)

type TaskPutDTO struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (h *Handler) handlePutTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskPutDTO
	err := json.NewDecoder(r.Body).Decode(&taskDTO)
	if err != nil {
		utils.RespondWithError(w, utils.ErrInvalidJson)
		return
	}
	taskId, err := strconv.ParseInt(taskDTO.ID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, utils.ErrGetTaskID)
		return
	}

	taskRequest := models.Task{
		ID:      taskId,
		Date:    taskDTO.Date,
		Title:   taskDTO.Title,
		Comment: taskDTO.Comment,
		Repeat:  taskDTO.Repeat,
	}

	task, err := validateTask(&taskRequest)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	err = h.repository.UpdateTask(task)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	utils.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
