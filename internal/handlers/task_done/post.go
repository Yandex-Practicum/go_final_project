package task_done

import (
	"net/http"
	"time"

	"go_final_project/internal/utils"
)

func (h *Handler) handlePostTaskDone(w http.ResponseWriter, r *http.Request) {
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

	if len(task.Repeat) == 0 {
		deleteErr := h.repository.DeleteTaskByID(task.ID)
		if deleteErr != nil {
			utils.RespondWithError(w, deleteErr)
			return
		}

		utils.SetJsonHeader(w)
		w.Write([]byte("{}"))
		return
	}

	task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	err = h.repository.UpdateTaskDate(task.ID, task.Date)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	utils.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
