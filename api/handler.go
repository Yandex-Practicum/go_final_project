package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/AlexJudin/go_final_project/usecases"
)

type errResponse struct {
	Error string `json:"error"`
}

type TaskHandler struct {
	uc usecases.Task
}

func NewTaskHandler(uc usecases.Task) TaskHandler {
	return TaskHandler{uc: uc}
}

func (h *TaskHandler) GetNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		log.Printf("Failed to parse time. Error: %+v", err)
		respErr := errResponse{
			Error: errors.New("Bad Request").Error(),
		}
		returnErr(http.StatusBadRequest, respErr, w)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := h.uc.GetNextDate(nowTime, date, repeat)
	if err != nil {
		log.Printf("Failed to get next date. Error: %+v", err)
		respErr := errResponse{
			Error: errors.New("Internal Server Error").Error(),
		}
		returnErr(http.StatusInternalServerError, respErr, w)
		return
	}

	nextDateJSON, err := json.Marshal(nextDate)
	if err != nil {
		log.Printf("Failed to marshal next date. Error: %+v", err)
		respErr := errResponse{
			Error: errors.New("Internal Server Error").Error(),
		}
		returnErr(http.StatusInternalServerError, respErr, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(nextDateJSON)
}

func returnErr(s int, e interface{}, w http.ResponseWriter) {
	jb, err := json.Marshal(e)
	if err != nil {
		s = http.StatusInternalServerError
		jb = []byte("{\"error\":\"" + "Internal Server Error" + "\"}")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	w.Write(jb)
}
