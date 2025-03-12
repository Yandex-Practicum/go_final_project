package handler

import (
	"encoding/json"
	"go_final_project/function"
	"go_final_project/structurs"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Repo HandlerRepository
}
type HandlerRepository interface {
	AddDB(DBFile string) error
	CheckTable() error
	CreateTable() error
	AddTask(t structurs.Task) (int, error)
	GetTasks() ([]structurs.Tasks, error)
	GetTaskId(id string) (structurs.Tasks, error)
	SearchTask(input string) ([]structurs.Tasks, error)
	PutTaskId(t structurs.Tasks) error
	DeleteTaskId(id string) error
	NextDate(now string, data string, repeat string) (string, error)
}

func NewHandler(repo HandlerRepository) Handler {
	return Handler{repo}
}
func (h Handler) PostGetPutDeleteTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.PostTask(w, r)
	case http.MethodGet:
		h.GetTaskId(w, r)
	case http.MethodDelete:
		h.DeleteTaskID(w, r)
	case http.MethodPut:
		h.PutTask(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
}
func (h Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	var TaskAdd structurs.Task
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&TaskAdd)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	repeat, err := function.RepeatChek(TaskAdd.Repeat)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	date, err := function.DataCheck(TaskAdd.Date)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	TaskAdd.Repeat = repeat
	TaskAdd.Date = date
	if TaskAdd.Title != "" && TaskAdd.Date != "" && TaskAdd.Repeat != "0" {
		id, errAdd := h.Repo.AddTask(TaskAdd)
		respId := strconv.Itoa(id)
		resp["id"] = respId
		if errAdd != nil {
			resp["error"] = errAdd.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
	} else {
		resp["error"] = "Не указан заголовок задачи"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
}
func (h Handler) GetTasksSearch(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	search := r.URL.Query().Get("search")
	if search != "" {
		input := function.SearcCheck(search)
		searchData, err := h.Repo.SearchTask(input)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		respSearch := make(map[string][]structurs.Tasks)
		respSearch["tasks"] = searchData
		err = json.NewEncoder(w).Encode(respSearch)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
	} else {
		tasks, err := h.Repo.GetTasks()
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		if tasks != nil {
			respTasks := make(map[string][]structurs.Tasks)
			respTasks["tasks"] = tasks
			err = json.NewEncoder(w).Encode(respTasks)
			if err != nil {
				resp["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			}

		} else {
			nulSlais := make([]string, 0)
			respNil := make(map[string][]string)
			respNil["tasks"] = nulSlais
			err := json.NewEncoder(w).Encode(respNil)
			if err != nil {
				resp["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}
}

func (h Handler) GetTaskId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := h.Repo.GetTaskId(id)
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		resp["error"] = "Задача не найдена"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

func (h Handler) PutTask(w http.ResponseWriter, r *http.Request) {
	var TaskChange structurs.Tasks
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewDecoder(r.Body).Decode(&TaskChange)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	id, err := function.IdCheck(TaskChange.Id)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	date, err := function.DataCheck(TaskChange.Date)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	repeat, err := function.RepeatChek(TaskChange.Repeat)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	TaskChange.Id = id
	TaskChange.Repeat = repeat
	TaskChange.Date = date
	if TaskChange.Title != "" && TaskChange.Date != "" && TaskChange.Repeat != "0" {
		err = h.Repo.PutTaskId(TaskChange)
		if err != nil {
			resp["error"] = "Задача не найдена"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		out := &structurs.Empty{}
		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
	} else {
		resp["error"] = "Задача не найдена"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

func (h Handler) DoneTaskId(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	nowNow := time.Now()
	format := "20060102"
	now := nowNow.Format(format)
	id := r.URL.Query().Get("id")
	task, err := h.Repo.GetTaskId(id)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	date := task.Date
	repeat := task.Repeat
	if task.Repeat == "" {
		err := h.Repo.DeleteTaskId(id)
		if err != nil {
			resp["error"] = "Не завершилась"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			out := &structurs.Empty{}
			err = json.NewEncoder(w).Encode(out)
			if err != nil {
				resp["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	} else {
		data, err := h.Repo.NextDate(now, date, repeat)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		task.Date = data
		err = h.Repo.PutTaskId(task)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			out := &structurs.Empty{}
			err = json.NewEncoder(w).Encode(out)
			if err != nil {
				resp["error"] = err.Error()
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}
	}
}
func (h Handler) DeleteTaskID(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := r.URL.Query().Get("id")
	idCheck, err := function.IdCheck(id)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = h.Repo.DeleteTaskId(idCheck)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		out := &structurs.Empty{}
		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
}
func (h Handler) NextData(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	nextdate, err := h.Repo.NextDate(now, date, repeat)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	answer, err := strconv.Atoi(nextdate)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	err = json.NewEncoder(w).Encode(answer)
	if err != nil {
		resp["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
}
