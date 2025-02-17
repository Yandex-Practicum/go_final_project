package handler

import (
	"encoding/json"
	"go_final_project/internal/model"
	"go_final_project/internal/service"
	"net/http"
	"strconv"
	"time"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) HandleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTask(w, r)
	case http.MethodPost:
		h.AddTask(w, r)
	case http.MethodPut:
		h.UpdateTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		http.Error(w, `{"error": "Метод не разрешен"}`, http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error": "Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	id, err := h.service.AddTask(task)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]interface{}{"id": id}
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error": "Ошибка декодирования JSON"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateTask(task); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат идентификатора задачи"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(id); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	limit := 50

	var tasks []model.Task
	var err error

	if search != "" {
		var date time.Time
		if date, err = time.Parse("02.01.2006", search); err == nil {
			dateStr := date.Format("20060102")
			tasks, err = h.service.GetTasksByDate(dateStr, limit)
		} else {
			tasks, err = h.service.GetTasksBySearch(search, limit)
		}
	} else {
		tasks, err = h.service.GetAllTasks(limit)
	}

	if err != nil {
		http.Error(w, `{"error": "Ошибка при получении задач"}`, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	response := map[string][]model.Task{"tasks": tasks}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) TaskDone(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат идентификатора задачи"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.TaskDone(id); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
