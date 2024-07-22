package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type errResponse struct {
	Error string `json:"error"`
}

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
// @Param date body string true "Дата задачи"
// @Param title body string true "Заголовок задачи"
// @Param comment body string true "Комментарий к задаче"
// @Param repeat body string true "Правило повторения"
// @Success 201 {object} model.TaskResp
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/task [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var (
		task model.Task
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusBadRequest, errResp, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusBadRequest, errResp, w)
		return
	}

	err = checkTaskRequest(&task)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusBadRequest, errResp, w)
		return
	}

	taskResp, err := h.uc.CreateTask(&task)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusInternalServerError, errResp, w)
		return
	}

	resp, err := json.Marshal(taskResp)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusInternalServerError, errResp, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// GetTasks ... Получить список ближайших задач
// @Summary Получить список ближайших задач
// @Description Получить список ближайших задач
// @Security ApiKeyAuth
// @Accept json
// @Tags Task
// @Success 200 {object} model.TasksResp
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/tasks [get]
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasksResp, err := h.uc.GetTasks()
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusInternalServerError, errResp, w)
		return
	}

	resp, err := json.Marshal(tasksResp)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusInternalServerError, errResp, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// GetTask ... Получить задачу
// @Summary Получить задачу
// @Description Получить задачу
// @Security ApiKeyAuth
// @Accept json
// @Tags Task
// @Param id query string true "Идентификатор задачи"
// @Success 200 {object} model.TaskResp
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/task [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	if taskId == "" {
		errResp := errResponse{
			Error: fmt.Errorf("task id is empty").Error(),
		}
		returnErr(http.StatusBadRequest, errResp, w)
		return
	}

	taskResp, err := h.uc.GetTaskById(taskId)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusInternalServerError, errResp, w)
		return
	}

	resp, err := json.Marshal(taskResp)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusInternalServerError, errResp, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// UpdateTask ... Редактировать задачу
// @Summary Редактировать задачу
// @Description Редактировать задачу
// @Security ApiKeyAuth
// @Accept json
// @Tags Task
// @Param id body string true "Идентификатор задачи"
// @Param date body string true "Дата задачи"
// @Param title body string true "Заголовок задачи"
// @Param comment body string true "Комментарий к задаче"
// @Param repeat body string true "Правило повторения"
// @Success 200 {string} string
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Router /api/task [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var (
		task model.Task
		buf  bytes.Buffer
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusBadRequest, errResp, w)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusBadRequest, errResp, w)
		return
	}

	err = checkTaskRequest(&task)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusBadRequest, errResp, w)
		return
	}

	err = h.uc.UpdateTask(&task)
	if err != nil {
		errResp := errResponse{
			Error: err.Error(),
		}
		returnErr(http.StatusInternalServerError, errResp, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(""))
}

func checkTaskRequest(task *model.Task) error {
	if task.Title == "" {
		return fmt.Errorf("task title is empty")
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	return nil
}

func returnErr(status int, message interface{}, w http.ResponseWriter) {
	messageJson, err := json.Marshal(message)
	if err != nil {
		status = http.StatusInternalServerError
		messageJson = []byte("{\"error\":\"" + err.Error() + "\"}")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(messageJson)
}
