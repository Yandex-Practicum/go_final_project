package task

import (
	"net/http"

	"go_final_project/internal/utils"
)

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIDFromQuery(r)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	err = h.repository.DeleteTaskByID(id)
	if err != nil {
		utils.RespondWithError(w, err)
		return
	}

	utils.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
