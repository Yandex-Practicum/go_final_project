package task

import (
	"net/http"
	"strconv"

	"go_final_project/internal/utils"
)

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
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

	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	_, deleteErr := h.db.Exec(deleteQuery, id)
	if deleteErr != nil {
		utils.RespondWithError(w, err.Error())
		return
	}

	utils.SetJsonHeader(w)
	w.Write([]byte("{}"))
}
