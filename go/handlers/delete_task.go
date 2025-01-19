package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var id int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.DeleteTask(id)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := json.Marshal(struct{}{})
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
