package deletetask

import (
	"encoding/json"
	"final_project/internal/common"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.FormValue("id")
	if _, err := strconv.Atoi(id); err != nil || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "неверно задан идентификатор задачи"})
		return
	}
	err := h.rep.DeleteTask(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Response{Error: "не удалось удалить задачу" + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(common.Response{})
}
