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
		utils.RespondWithError(w, err)
		return
	}

	task.ID, err = h.repository.CreateTask(task)
	if err != nil {
		utils.RespondWithError(w, err)
	}

	log.Printf("Задача добавлена: %+v\n", task)

	response := models.Response{ID: task.ID}
	utils.Respond(w, response)
}
