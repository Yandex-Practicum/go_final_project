package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AlexJudin/go_final_project/usecases"
	"github.com/AlexJudin/go_final_project/usecases/model"
)

type TaskHandler struct {
	uc usecases.Task
}

func NewTaskHandler(uc usecases.Task) TaskHandler {
	return TaskHandler{uc: uc}
}

// GetNextDate ... Получение следующей даты
// @Summary Получение следующей даты
// @Description Получение следующей даты
// @Security ApiKeyAuth
// @Accept json
// @Tags nextDate
// @Success 200 {string} []byte
// @Failure 400,503 {object} http.Error
// @Router /api/nextdate [get]
func (h *TaskHandler) GetNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		log.Printf("Failed to parse time. Error: %+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := h.uc.GetNextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("Failed to get next date. Error: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// CreateTask ... Добавить новую задачу
// @Summary Добавить новую задачу
// @Description Добавить новую задачу
// @Security ApiKeyAuth
// @Accept json
// @Tags Task
// @Param name body string true "Наименование контрагента" maxlength(300)
// @Success 201 {object} model.TaskResp
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /api/task [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var (
		task model.TaskReq
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "title is empty", http.StatusBadRequest)
		return
	}

	taskResp, err := h.uc.CreateTask(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(taskResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
