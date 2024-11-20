package handlers

import (
	"encoding/json"
	"final_project/internal/common"

	"net/http"
)

func (h *Handler) EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	task := common.Task{}
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "сервер не может обработать данные отправленные вами"})
	}
	switch "" {
	case task.ID:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "не указан идентификатор"})
		return
	case task.Date:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "не указана дата"})
		return
	case task.Title:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "не указан заголовок"})
		return
	case task.Comment:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "не указан комментарий"})
		return
	case task.Repeat:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.Response{Error: "не указано правило повторения"})
		return
	}

	err = h.rep.EditTask(task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(common.Response{Error: "сервер не может изменить данные " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(common.Response{})
}
