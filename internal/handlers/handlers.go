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

func (h *TaskHandler) CRUDTask(write http.ResponseWriter, request *http.Request) {

	now := time.Now().Format(constants.DateFormat)
	newTask := new(models.Task)

	switch request.Method {

	case http.MethodGet:

		id := request.FormValue("id")

		task, err := h.storage.GetOneTask(id)
		if err != nil {
			myfunctions.WriteJSONError(write, http.StatusBadRequest, fmt.Sprintf("GetOneTask: function error: %v", err))
			return
		}

		write.Header().Set("Content-Type", "application/json")
		write.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(write).Encode(task); err != nil {
			http.Error(write, "failed to encode response", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:

		if err := json.NewDecoder(request.Body).Decode(newTask); err != nil {
			http.Error(write, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		if err := myfunctions.ValidateTaskRequest(newTask, now); err != nil {
			myfunctions.WriteJSONError(write, http.StatusBadRequest, fmt.Sprintf("ValidateTaskRequest: function error: %v", err))
			return
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

	case http.MethodPut:

		if err := json.NewDecoder(request.Body).Decode(newTask); err != nil {
			http.Error(write, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		if err := myfunctions.ValidateTaskRequest(newTask, now); err != nil {
			myfunctions.WriteJSONError(write, http.StatusBadRequest, fmt.Sprintf("ValidateTaskRequest: function error: %v", err))
			return
		}

		if err := h.storage.EditTask(*newTask); err != nil {
			myfunctions.WriteJSONError(write, http.StatusInternalServerError, fmt.Sprintf("EditTask: function error: %v", err))
			return
		}

		write.Header().Set("Content-Type", "application/json")
		write.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(write).Encode(map[string]interface{}{}); err != nil {
			http.Error(write, "failed to encode response", http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:

		id := request.FormValue("id")

		if err := h.storage.DeleteTask(id); err != nil {
			myfunctions.WriteJSONError(write, http.StatusInternalServerError, fmt.Sprintf("DeleteTask: function error: %v", err))
			return
		}

		write.Header().Set("Content-Type", "application/json")
		write.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(write).Encode(struct{}{}); err != nil {
			http.Error(write, "failed to encode response", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(write, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *TaskHandler) GetTasks(write http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(write, "invalid method", http.StatusMethodNotAllowed)
	}

	searchQuery := request.FormValue("search")

	if searchQuery != "" {
		searchTasks, err := h.storage.SearchTasks(searchQuery)
		if err != nil {
			myfunctions.WriteJSONError(write, http.StatusInternalServerError, fmt.Sprintf("SearchTasks: function error: %v", err))
			return
		}

		write.Header().Set("Content-Type", "application/json")
		write.WriteHeader(http.StatusOK)

		response := map[string][]models.Task{"tasks": searchTasks}

		if err := json.NewEncoder(write).Encode(response); err != nil {
			http.Error(write, "failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	tasks, err := h.storage.GetTasks()
	if err != nil {
		http.Error(write, fmt.Sprintf("GetTasks: function error: %v", err), http.StatusInternalServerError)
		return
	}

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)

	response := map[string][]models.Task{"tasks": tasks}

	if err := json.NewEncoder(write).Encode(response); err != nil {
		http.Error(write, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *TaskHandler) TaskIsDone(write http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		http.Error(write, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := request.FormValue("id")

	task, err := h.storage.GetOneTask(id)
	if err != nil {
		myfunctions.WriteJSONError(write, http.StatusBadRequest, fmt.Sprintf("GetOneTask: function error: %v", err))
		return
	}

	if task.Repeat == "" {
		if err := h.storage.DeleteTask(id); err != nil {
			myfunctions.WriteJSONError(write, http.StatusInternalServerError, fmt.Sprintf("DeleteTask: function error: %v", err))
			return
		}
	} else {
		nextDate, err := myfunctions.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			myfunctions.WriteJSONError(write, http.StatusInternalServerError, fmt.Sprintf("NextDate error: %v", err))
			return
		}

		task.Date = nextDate

		if err := h.storage.EditTask(task); err != nil {
			myfunctions.WriteJSONError(write, http.StatusInternalServerError, fmt.Sprintf("EditTask error: %v", err))
			return
		}
	}

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(write).Encode(struct{}{}); err != nil {
		http.Error(write, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) Authentication(write http.ResponseWriter, request *http.Request) {

	type authResponse struct {
		Token string `json:"token"`
	}

	type password struct {
		Password string `json:"password"`
	}

	pwdFromJSON := password{
		Password: "",
	}

	if request.Method != http.MethodPost {
		http.Error(write, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(request.Body).Decode(&pwdFromJSON); err != nil {
		http.Error(write, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	pwd := pwdFromJSON.Password

	if pwd == "" {
		myfunctions.WriteJSONError(write, http.StatusBadRequest, "password cannot be empty")
		return
	}

	token, err := myfunctions.PwdValidateGenerateJWT(pwd)
	if err != nil {
		myfunctions.WriteJSONError(write, http.StatusBadRequest, fmt.Sprintf("PwdValidateGenerateJWT: function error: %v", err))
		return
	}

	response := authResponse{Token: token}

	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(write).Encode(response); err != nil {
		http.Error(write, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
