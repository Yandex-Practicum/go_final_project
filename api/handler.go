package api

import (
	"net/http"
	"time"

	"github.com/AlexJudin/go_final_project/usecases"
)

type TaskHandler struct {
	uc usecases.Task
}

func NewTaskHandler(uc usecases.Task) TaskHandler {
	return TaskHandler{uc: uc}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		return
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	answer, err := h.uc.NextDate(nowTime, date, repeat)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}
