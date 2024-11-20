package handlers

import (
	"encoding/json"
	"final_project/internal/common"

	"net/http"
	"strconv"
)

func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	w.Header().Set("Content-Type", "application/json")
	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "указанный идентификатор не целочисленного значения или пуст"})
		return
	}
	task, err := h.rep.GetTaskByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Response{Error: err.Error()})
		return
	}
	if task.Date == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Response{Error: "задачи с таким id не сущетсвует"})
		return
	}
	json.NewEncoder(w).Encode(task)
}
