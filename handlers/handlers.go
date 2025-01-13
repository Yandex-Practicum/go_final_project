package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/handlers" // Correct
	"go_final_project/storage"  // Correct
	"go_final_project/tests"    //Correct
)

const ParseDate = "20060102"

type Handlers struct {
	TaskStorage storage.DB
}

func (h *Handlers) respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": message})
}

func (h *Handlers) AddTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task task.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := task.Checktitle(); err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		taskmod, err := task.Checkdate()
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err = taskmod.Countdate(); err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		id, err := h.TaskStorage.Addtasktodb(taskmod)
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

func (h *Handlers) ChangeTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task task.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			h.respondWithError(w, http.StatusBadRequest, "Deserialization error JSON")
			return
		}
		if errstr := task.CheckId(); errstr != "" {
			h.respondWithError(w, http.StatusBadRequest, errstr)
			return
		}
		if err := task.Checktitle(); err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		_, err := task.Checkdate()
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errstr := task.CheckRepeate(); errstr != "" {
			h.respondWithError(w, http.StatusBadRequest, errstr)
			return
		}
		if errstr := h.TaskStorage.Update(task); errstr != "" {
			h.respondWithError(w, http.StatusInternalServerError, errstr)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func (h *Handlers) GetTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			h.respondWithError(w, http.StatusBadRequest, "No identifier specified")
			return
		}
		task, err := h.TaskStorage.Findtask(id)
		if err != "" {
			h.respondWithError(w, http.StatusInternalServerError, err)
			return
		}
		json.NewEncoder(w).Encode(task)
	}
}

func (h *Handlers) TaskDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			h.respondWithError(w, http.StatusBadRequest, "No identifier specified")
			return
		}
		task, err := h.TaskStorage.Findtask(id)
		if err != "" {
			h.respondWithError(w, http.StatusInternalServerError, "Task not found")
			return
		}
		if task.Repeat == "" {
			// Удаляем одноразовую задачу
			if err = h.TaskStorage.DeleteQuery(id); err != "" {
				h.respondWithError(w, http.StatusInternalServerError, "Error deleting task")
				return
			}
		} else {
			// Рассчитываем следующую дату для периодической задачи
			now := time.Now()
			nextDate, err := nextdate.CalcNextDate(now.Format(ParseDate), task.Date, task.Repeat)
			if err != nil {
				h.respondWithError(w, http.StatusInternalServerError, "Error calculating next date")
				return
			}

			// Обновляем дату задачи
			if err = h.TaskStorage.Updatetask(nextDate, id); err != "" {
				h.respondWithError(w, http.StatusInternalServerError, "Task update error")
				return
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func (h *Handlers) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			h.respondWithError(w, http.StatusBadRequest, "No identifier specified")
			return
		}
		if err := h.TaskStorage.DeleteQuery(id); err != "" {
			h.respondWithError(w, http.StatusInternalServerError, "Error deleting task")
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

func (h *Handlers) NextDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := r.FormValue("now")
		date := r.FormValue("date")
		repeat := r.FormValue("repeat")
		if repeat == "" || now == "" || date == "" {
			h.respondWithError(w, http.StatusBadRequest, "Incorrect data specified in the request")
			return
		}
		nextDate, err := nextdate.CalcNextDate(now, date, repeat)
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, "Incorrect data specified in the request")
			return
		}
		w.Write([]byte(nextDate))
	}
}

func (h *Handlers) ReceiveTasks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := h.TaskStorage.GetTasks()
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		response := map[string]interface{}{
			"tasks": tasks,
		}
		json.NewEncoder(w).Encode(response)
	}
}
