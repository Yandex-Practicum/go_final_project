package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"todo_restapi/internal/dto"
	"todo_restapi/internal/myfunctions"
	"todo_restapi/internal/storage"
	"todo_restapi/pkg/constants"
)

type TaskHandler struct {
	storage *storage.Storage
}

func NewTaskHandler(storage *storage.Storage) *TaskHandler {
	return &TaskHandler{storage: storage}
}

func NextDateHandler(write http.ResponseWriter, request *http.Request) {

	timeNow := request.FormValue("now")
	date := request.FormValue("date")
	repeat := request.FormValue("repeat")

	timeParse, err := time.Parse(constants.DateFormat, timeNow)
	if err != nil {
		http.Error(write, fmt.Sprintf("time parse error: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := myfunctions.NextDate(timeParse, date, repeat)
	if err != nil {
		http.Error(write, fmt.Sprintf("NextDate: function error: %v", err), http.StatusBadRequest)
		return
	}

	write.WriteHeader(http.StatusOK)

	if _, err := write.Write([]byte(result)); err != nil {
		http.Error(write, "failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) AddTask(write http.ResponseWriter, request *http.Request) {

	now := time.Now().Format(constants.DateFormat)
	newTask := new(models.Task)

	switch request.Method {

	case http.MethodGet:
	case http.MethodPost:

		if err := json.NewDecoder(request.Body).Decode(newTask); err != nil {
			http.Error(write, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		if newTask.Title == "" {
			myfunctions.WriteJSONError(write, http.StatusBadRequest, "title is empty")
			return
		}

		if newTask.Date == "" || newTask.Date == "today" {
			newTask.Date = now
		}

		newDate, err := time.Parse(constants.DateFormat, newTask.Date)
		if err != nil {
			myfunctions.WriteJSONError(write, http.StatusBadRequest, "invalid date format")
			return
		}

		if newDate.Before(time.Now()) {
			if newTask.Repeat == "" {
				newTask.Date = now
			} else {
				dateCalculation, err := myfunctions.NextDate(time.Now(), newTask.Date, newTask.Repeat)
				if err != nil {
					myfunctions.WriteJSONError(write, http.StatusBadRequest, fmt.Sprintf("NextDate: function error: %v", err))
					return
				}
				newTask.Date = dateCalculation
			}
		}

		taskID, err := h.storage.AddTask(*newTask)
		if err != nil {
			http.Error(write, fmt.Sprintf("AddTask: add task error: %v", err), http.StatusInternalServerError)
			return
		}

		write.Header().Set("Content-Type", "application/json")
		write.WriteHeader(http.StatusCreated)

		response := map[string]int64{"id": taskID}

		if err := json.NewEncoder(write).Encode(response); err != nil {
			http.Error(write, "failed to encode response", http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
	default:
		http.Error(write, "method not allowed", http.StatusMethodNotAllowed)
	}

}
