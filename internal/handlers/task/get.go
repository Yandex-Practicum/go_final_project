package task

import (
	"net/http"
	"strconv"

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
	id, err := utils.GetIDFromQuery(r)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	task, err := h.repository.GetTaskByID(id)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	response := GetTaskResponseDTO{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	utils.Respond(w, response)
}
