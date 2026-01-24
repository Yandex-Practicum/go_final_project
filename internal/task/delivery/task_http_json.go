package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test/internal/task"
)

type ServiceInterface interface {
	NextDate(now string, date string, repeat string) (error, string)
	Create(t *task.Task) (string, error)
	GetAll(search string) (*task.List, error)
	GetById(id int) (*task.Task, error)
	Update(t *task.Task) (*task.Task, error)
	Delete(id int) error
	Done(paramId string) error
}

type TaskHttp struct {
	service ServiceInterface
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewTaskHttp(service ServiceInterface) *TaskHttp {
	return &TaskHttp{
		service: service,
	}
}

func (th *TaskHttp) HandleTime(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	err, newDate := th.service.NextDate(now, date, repeat)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	response, _ := strconv.Atoi(newDate)
	writeResponse(response, w, false)
}

func (th *TaskHttp) Create(w http.ResponseWriter, r *http.Request) {
	var t task.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		errorResponse := ErrorResponse{Error: "Ошибка десериализации JSON"}
		writeResponse(errorResponse, w, true)
		return
	}

	id, err := th.service.Create(&t)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	response := task.Task{ID: id}
	writeResponse(response, w, false)
}

func (th *TaskHttp) GetList(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	taskList, err := th.service.GetAll(search)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	writeResponse(taskList, w, false)
}

func (th *TaskHttp) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	t, err := th.service.GetById(id)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	writeResponse(t, w, false)
}

func (th *TaskHttp) Delete(w http.ResponseWriter, r *http.Request) {
	paramId := r.FormValue("id")
	if paramId == "" {
		errorResponse := ErrorResponse{Error: "задача не найдена"}
		writeResponse(errorResponse, w, true)
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	err = th.service.Delete(id)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	writeResponse(task.Task{}, w, false)
}

func (th *TaskHttp) Edit(w http.ResponseWriter, r *http.Request) {
	var t task.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		errorResponse := ErrorResponse{Error: "Ошибка десериализации JSON"}
		writeResponse(errorResponse, w, true)
		return
	}

	updatedTask, err := th.service.Update(&t)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	writeResponse(updatedTask, w, false)
}

func (th *TaskHttp) Done(w http.ResponseWriter, r *http.Request) {
	paramId := r.FormValue("id")
	err := th.service.Done(paramId)
	if err != nil {
		errorResponse := ErrorResponse{Error: err.Error()}
		writeResponse(errorResponse, w, true)
		return
	}

	writeResponse(task.Task{}, w, false)
}

func writeResponse(data interface{}, w http.ResponseWriter, isError bool) {
	toJson, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, `Marshal error`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(toJson)
	if err != nil {
		http.Error(w, `Failed response`, http.StatusInternalServerError)
		return
	}

	if !isError {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
